# Rationale

markdown-style enforces consistent formatting for all
generated markdown files in the project.

## Why 72–80 character wrap instead of sentence-per-line

Sentence-per-line produces enormous diffs when
editing mid-paragraph. Fixed-width wrap keeps diffs
localized to the changed lines. It also renders
predictably in terminals and editors without soft wrap.

## Why only dash is allowed for unordered lists

Eliminating marker choice (*, +, -) removes a
decision agents would make inconsistently. One marker
means zero formatting drift across files.

## Why only ATX headings are allowed

Setext headings (underlines) break visually when the
heading text is longer than the underline. ATX (#)
scales to any length and any heading level without
ambiguity.

## Why code blocks and links are exempt from line length

Breaking URLs mid-line makes them unclickable. Breaking
code blocks changes semantics. These elements must be
exempt or agents produce broken output trying to hit
the 80-char target.
