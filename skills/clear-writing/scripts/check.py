#!/usr/bin/env python3
"""Flag sentences that are hard to read in Markdown or plain text.

Usage:
    check.py [--max-words N] [--changed REF] PATH...
    check.py - < text.md

PATH can be a file or a directory; directories are searched for
.md, .markdown, .mdx, .txt and .templ files. Code blocks, inline
code, frontmatter, headings, tables and HTML are skipped. In .templ
files only the text between tags inside components is read.

Prints one line per flagged sentence:

    path:line: reason, reason | sentence

Reasons marked "(weak)" are fine on their own; act on them only
when they make the sentence harder to read.

Exit status: 0 nothing flagged, 1 sentences flagged, 2 error.
"""

import argparse
import html
import os
import re
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
WATCH_LIST = os.path.join(HERE, "..", "references", "watch-list.md")
TEXT_EXTENSIONS = (".md", ".markdown", ".mdx", ".txt", ".templ")
SKIP_DIRS = {"node_modules", "vendor", "dist", "build"}

CODE = "‹code›"
ABBREVIATIONS = {"e.g.", "i.e.", "etc.", "vs.", "cf.", "approx.",
                 "no.", "mr.", "ms.", "dr.", "fig."}

LIST_ITEM = re.compile(r"^\s*(?:[-*+]|\d+[.)])\s+")
FENCE = re.compile(r"^\s*(`{3,}|~{3,})")
LINK_DEF = re.compile(r"^\s*\[[^\]]+\]:\s")
WORD = re.compile(r"[A-Za-z0-9][\w'’-]*|" + CODE)
BOUNDARY = re.compile(r"(?<=[.!?])[\"'’”)\]]*\s+(?=[A-Z0-9\"'“(\[‹])")
CONNECTORS = re.compile(
    r"\b(?:which|because|so that|while|whereas)\b|,\s+(?:and|but|so|or)\b",
    re.IGNORECASE)
DASH = re.compile(r"\s[—–]\s|\s--\s|\w—\w")
ASIDE = re.compile(r"\(([^()]*)\)")

TEMPL_START = re.compile(r"^templ\s+\w")
TEMPL_CONTROL = re.compile(
    r"^(?:(?:if|for|switch)\b.*\{|(?:case|default)\b.*:|[}@]|//)")
TAG_NAME = re.compile(r"<(/?)([A-Za-z][\w-]*)")
INLINE_TAGS = {"a", "abbr", "b", "bdi", "bdo", "br", "cite", "data", "dfn",
               "em", "i", "kbd", "mark", "q", "s", "samp", "small", "span",
               "strong", "sub", "sup", "time", "u", "var", "wbr"}
SKIPPED_TAGS = {"h1", "h2", "h3", "h4", "h5", "h6", "pre", "script",
                "style", "td", "th", "title", "textarea", "svg"}


def load_watch_list(path):
    """Return (label, regex) pairs for every term in the watch list."""
    terms = []
    section = None
    with open(path, encoding="utf-8") as f:
        for line in f:
            if line.startswith("## "):
                section = line[3:].strip().lower()
                continue
            if (section not in ("words", "phrases")
                    or not line.startswith("- ")):
                continue
            head = line.split("→")[0]
            for term in re.findall(r"`([^`]+)`", head):
                body = r"\s+".join(re.escape(w) for w in term.split())
                body = body.replace("'", "['’]")
                if section == "words":
                    body += r"(?:s|es|d|ed|ing|ly)?"
                pattern = re.compile(r"(?<![\w-])" + body + r"(?![\w-])",
                                     re.IGNORECASE)
                terms.append((len(term), f"{section[:-1]} '{term}'", pattern))
    # Longest first, so "a testament to" is reported instead of
    # "testament" inside it.
    terms.sort(key=lambda t: -t[0])
    return [(label, pattern) for _, label, pattern in terms]


def clean_inline(text):
    """Turn a Markdown line into plain prose for counting."""
    text = re.sub(r"`[^`]*`", CODE, text)
    text = re.sub(r"!\[([^\]]*)\]\([^)]*\)", r"\1", text)
    text = re.sub(r"\[([^\]]+)\]\([^)]*\)", r"\1", text)
    text = re.sub(r"\[([^\]]+)\]\[[^\]]*\]", r"\1", text)
    text = re.sub(r"<https?://[^>]+>|https?://\S+", "link", text)
    text = re.sub(r"<[^>]+>", "", text)
    return text.replace("**", "").replace("__", "")


def prose_blocks(lines):
    """Yield blocks of prose as lists of (line number, text)."""
    block = []
    fence = None
    in_frontmatter = bool(lines) and lines[0].strip() == "---"
    for number, raw in enumerate(lines, 1):
        line = raw.rstrip("\n")
        stripped = line.strip()
        if in_frontmatter:
            if number > 1 and stripped in ("---", "..."):
                in_frontmatter = False
            continue
        opening = FENCE.match(line)
        if fence:
            if opening and opening.group(1)[0] == fence[0]:
                fence = None
            continue
        if opening:
            fence = opening.group(1)
            yield block
            block = []
            continue
        if (not stripped or stripped.startswith(("#", "|", "<"))
                or LINK_DEF.match(line)):
            yield block
            block = []
            continue
        text = re.sub(r"^(?:>\s?)+", "", stripped)
        if LIST_ITEM.match(line):
            yield block
            block = []
            text = LIST_ITEM.sub("", line)
        block.append((number, clean_inline(text.strip())))
    yield block


def add_text(block, number, parts):
    """Add the text found on one templ line to a block."""
    text = re.sub(r"\s+", " ", html.unescape("".join(parts))).strip()
    if text:
        block.append((number, text))


def templ_blocks(lines):
    """Yield the page text of templ components as prose blocks.

    Only text between tags inside `templ Name() { ... }` is read. Go
    code, control flow, attributes, headings, tables, <pre>, <script>
    and <style> are skipped; `{ expr }` and <code> count as code.
    """
    block = []
    in_templ = False
    skip = None      # tag whose content is skipped until it closes
    tag = None       # a tag that continues on the next line
    comment = False
    expr = 0         # depth of a { expression } that spans lines
    for number, raw in enumerate(lines, 1):
        line = raw.rstrip("\n")
        if not in_templ:
            in_templ = TEMPL_START.match(line) is not None
            continue
        # A multi-line string in a tag or expression can hold any line,
        # so code and the component's end are only seen outside them.
        open_span = tag is not None or skip is not None or comment or expr
        if line.startswith("}") and not open_span:
            in_templ = False
            yield block
            block = []
            continue
        if not open_span and TEMPL_CONTROL.match(line.strip()):
            yield block
            block = []
            continue
        text = []
        i = 0
        while i < len(line):
            if comment:
                end = line.find("-->", i)
                if end < 0:
                    break
                comment, i = False, end + 3
            elif skip is not None:
                end = re.compile(r"</" + skip + r"\s*>", re.I).search(line, i)
                if end is None:
                    break
                skip, i = None, end.end()
            elif tag is not None:
                end = line.find(">", i)
                if end < 0:
                    tag += line[i:] + " "
                    break
                tag_text, tag = tag + line[i:end + 1], None
                i = end + 1
                closing, name = TAG_NAME.match(tag_text).groups()
                name = name.lower()
                if name in INLINE_TAGS:
                    text.append(" " if name == "br" else "")
                elif name == "code" and not closing:
                    text.append(CODE)
                    skip = "code"
                else:
                    if (name in SKIPPED_TAGS and not closing
                            and not tag_text.endswith("/>")):
                        skip = name
                    add_text(block, number, text)
                    yield block
                    block, text = [], []
            elif line.startswith("<!--", i):
                comment, i = True, i + 4
            elif re.match(r"</?[A-Za-z]", line[i:i + 3]):
                tag = ""
            elif expr or line[i] == "{":
                if not expr:
                    text.append(CODE)
                while i < len(line):
                    expr += {"{": 1, "}": -1}.get(line[i], 0)
                    i += 1
                    if expr == 0:
                        break
            else:
                end = min([k for k in (line.find("<", i), line.find("{", i))
                           if k > i] or [len(line)])
                text.append(line[i:end])
                i = end
        add_text(block, number, text)
    yield block


def sentences(block):
    """Split a block into (first line, last line, sentence) tuples."""
    text = ""
    starts = []
    for number, part in block:
        if text:
            text += " "
        starts.append((len(text), number))
        text += part

    def line_at(offset):
        found = starts[0][1]
        for start, number in starts:
            if start <= offset:
                found = number
        return found

    begin = 0
    for match in BOUNDARY.finditer(text):
        before = text[begin:match.start()].split()
        if before and before[-1].lower() in ABBREVIATIONS:
            continue
        yield line_at(begin), line_at(match.start()), text[begin:match.start()]
        begin = match.end()
    if text[begin:].strip():
        yield line_at(begin), line_at(len(text) - 1), text[begin:]


def reasons_for(sentence, max_words, terms):
    words = len(WORD.findall(sentence))
    reasons = []
    if words > max_words:
        reasons.append(f"long ({words} words)")
    clauses = len(CONNECTORS.findall(sentence))
    if clauses >= 3 or (clauses >= 2 and words >= 20):
        reasons.append(f"tangled ({clauses + 1} clauses)")
    taken = set()
    for label, pattern in terms:
        match = pattern.search(sentence)
        if match and taken.isdisjoint(range(*match.span())):
            taken.update(range(*match.span()))
            reasons.append(label)
    if DASH.search(sentence):
        reasons.append("dash joins clauses (weak)")
    if ";" in sentence.replace(CODE, ""):
        reasons.append("semicolon (weak)")
    for aside in ASIDE.findall(sentence):
        if len(WORD.findall(aside)) > 8:
            reasons.append("long aside in parentheses (weak)")
            break
    return reasons


def changed_lines(path, ref):
    """Return the line numbers of path that differ from ref in git."""
    folder, name = os.path.split(os.path.abspath(path))
    out = subprocess.run(
        ["git", "-C", folder, "diff", "-U0", ref, "--", name],
        capture_output=True, text=True, check=True).stdout
    lines = set()
    for start, count in re.findall(r"^@@ .* \+(\d+)(?:,(\d+))? @@", out, re.M):
        first = int(start)
        lines.update(range(first, first + int(count or "1")))
    return lines


def text_files(paths):
    for path in paths:
        if path == "-" or os.path.isfile(path):
            yield path
            continue
        if not os.path.isdir(path):
            raise FileNotFoundError(path)
        for root, dirs, files in os.walk(path):
            dirs[:] = sorted(d for d in dirs
                             if not d.startswith(".") and d not in SKIP_DIRS)
            for name in sorted(files):
                if name.lower().endswith(TEXT_EXTENSIONS):
                    yield os.path.join(root, name)


def main():
    parser = argparse.ArgumentParser(
        description="Flag sentences that are hard to read.")
    parser.add_argument("paths", nargs="+", metavar="PATH",
                        help="file, directory, or - for standard input")
    parser.add_argument("--max-words", type=int, default=25,
                        help="flag sentences longer than this (default 25)")
    parser.add_argument("--changed", metavar="REF",
                        help="only report sentences on lines changed since "
                             "the git REF")
    args = parser.parse_args()

    try:
        terms = load_watch_list(WATCH_LIST)
        flagged = 0
        for path in text_files(args.paths):
            if path == "-":
                name, lines = "<stdin>", sys.stdin.read().splitlines(True)
            else:
                with open(path, encoding="utf-8") as f:
                    name, lines = path, f.readlines()
            keep = None
            if args.changed and path != "-":
                keep = changed_lines(path, args.changed)
            if name.lower().endswith(".templ"):
                blocks = templ_blocks(lines)
            else:
                blocks = prose_blocks(lines)
            for block in blocks:
                if not block:
                    continue
                for first, last, sentence in sentences(block):
                    if keep is not None and not keep.intersection(
                            range(first, last + 1)):
                        continue
                    reasons = reasons_for(sentence, args.max_words, terms)
                    if reasons:
                        flagged += 1
                        print(f"{name}:{first}: {', '.join(reasons)} | "
                              f"{sentence.strip()}")
    except subprocess.CalledProcessError as err:
        print(f"check.py: git diff failed: {err.stderr.strip()}",
              file=sys.stderr)
        return 2
    except OSError as err:
        print(f"check.py: {err}", file=sys.stderr)
        return 2
    return 1 if flagged else 0


if __name__ == "__main__":
    sys.exit(main())
