---
name: clear-writing
description: >
  Finds sentences in Markdown, text and templ files that are long,
  tangled or sound machine-written, rewrites only those in place in
  plain English, and reports the flagged sentences it kept. Use when
  the user wants docs, a README or page copy easier to read, asks to
  simplify, shorten or tighten wordy sentences, to deslop or humanize
  text, to remove AI writing tells, or to check the writing in a
  branch before a pull request.
metadata:
  version: "1.0.0"
  tags: "category:documentation, category:review, language:markdown"
---

# Clear Writing

Make writing easy to read for every reader, including people who
read English as a second language. A script finds the sentences that
need work, you judge each one, and you rewrite only the ones that
need it. Every other sentence stays exactly as the author wrote it.

## Workflow

1. **Pick the scope.** Use the files, folder or pasted text the user
   names. For "the branch" or "before a PR", use the changed lines:
   `--changed <base>`, where `<base>` is the branch the work will
   merge into, such as `main`.

2. **Run the check.** The script is `scripts/check.py` in this
   skill's folder. Run it from the project root:

   ```bash
   python3 <skill-folder>/scripts/check.py [--changed main] <paths>
   ```

   In `.templ` files only the page text inside components is read,
   not the Go code or the attributes.

   Pasted text goes in on standard input with `-` as the path. Each
   output line is `path:line: reasons | sentence`. Exit status 1
   means sentences were flagged and 0 means none were.

   If Python 3 is not available, read the text yourself and flag
   these:
   - sentences of more than 25 words
   - sentences of three or more clauses
   - the terms in [the watch list](references/watch-list.md)
   - clauses joined by a dash, unless at most three words come
     before it
   - clauses joined by a semicolon (weak)
   - asides of more than eight words in parentheses (weak)

3. **Judge each flagged sentence.** Read it in its paragraph and
   decide whether it is really hard to read.
   - Fix every sentence flagged `long` or `tangled`.
   - Fix a watch-list term unless it is the precise word in this
     context, such as `robust` for a statistical method.
   - Fix a dash that joins clauses, unless the sentence is a slogan
     or tagline, or a quote.
   - Leave a `(weak)` reason alone unless it makes the sentence
     harder to read, or the sentence has other reasons too.

4. **Rewrite in place.** Edit the file directly. For each sentence
   you fix, apply only the fixes its reasons call for:
   - **Long or tangled:** split it so each sentence holds one idea.
     Put a condition before the instruction it limits. Turn three
     or more parallel items into a list.
   - **Watch-list term:** use the plain alternative from the watch
     list, or delete the term when the list says so.
   - **Dash or semicolon:** use a period, a comma or a colon.
   - **Long aside:** make it its own sentence.

   When you split a sentence, name the subject again instead of
   starting with "it" or "this" when the reader could lose track.
   Keep the facts, terms, tone and contractions of the original.
   Keep the file's line wrapping.

5. **Check again.** Run the script on the same scope. Fix any
   sentence your rewrite made long. A flagged sentence you decided
   to keep stays flagged.

6. **Report.** Tell the user how many sentences you rewrote. List
   each flagged sentence you kept, with the reason. Leave the
   changes for the user to review in the diff.

To run the check on every commit, set up the git hooks in
[references/git-hooks.md](references/git-hooks.md). Do this only
when the user asks for it.

## Example

Input `docs/skills.md`:

```markdown
Skills are context-efficient by design. Your agent reads a one-line
description at session start — the full workflow only loads when
it's actually triggered.
```

Check output:

```text
docs/skills.md:1: phrase 'by design' | Skills are context-efficient by design.
docs/skills.md:1: word 'actually', dash joins clauses | Your agent reads a one-line description at session start — the full workflow only loads when it's actually triggered.
```

Rewritten in place:

```markdown
Skills are context-efficient. Your agent reads a one-line
description at session start. It loads the full workflow only when
the skill is triggered.
```

`by design` and `actually` were deleted as the watch list says, and
the dash became a period.

## Boundaries

**DOES:**
- Find long, tangled and machine-sounding sentences with a script
- Rewrite flagged sentences in place, keeping meaning and tone
- Work on files, folders, pasted text or a branch's changed lines

**Does NOT:**
- Rewrite sentences the check did not flag
- Change code blocks, inline code, commands, paths, links,
  headings, frontmatter, tables or quoted text
- Fix spelling or grammar, translate, or change the document's
  structure
- Try to hide that text was written with AI

## Common Failures

- **Rewriting unflagged text.** The whole point is a small diff.
  NEVER "improve" a sentence the check did not flag.
- **Losing a fact when splitting.** Each number, name, condition
  and limit in the original MUST appear in the rewrite.
- **Telegraph style.** Short does not mean clipped. Keep articles,
  verbs and "that". "Run check. Fix output." is worse than the
  original.
- **Breaking Markdown.** Splitting across a link or inline code
  breaks it. Keep each link and code span whole.
- **Breaking markup.** In a templ file, change only the text between
  tags. Keep every tag, attribute and `{ expression }` as it was.
