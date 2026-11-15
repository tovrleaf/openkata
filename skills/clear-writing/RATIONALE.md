# Rationale

clear-writing finds the sentences that are hard to read and rewrites
only those, so the author's voice survives and the diff stays small.

## Why only flagged sentences are rewritten

Rewriting a whole document flattens it into one register, even where
the original was already clear. Fixing only flagged sentences keeps
what works, and gives a diff the author can review in minutes.

## Why a script finds them and the agent judges them

Agents count words badly and drift between runs. A script makes the
search exact and repeatable. But a long sentence can be a clear list,
and `robust` can be the right technical word. The agent decides which
hits to fix.

## Why 25 words

It is the limit ASD-STE100 Simplified Technical English sets for
descriptive text, which has been tested on readers of English as a
second language. A stricter limit splits sentences that read fine.

## Why weak signs are acted on only together

A dash, a semicolon or an aside is normal punctuation. On its own it
says little, and banning it would rewrite good sentences. Next to a
long clause or an inflated word, it is part of the problem.

## Why it edits in place

The rewrite only matters once it is in the file. Editing in place
and showing the diff is faster than proposing every change, and
version control already makes each change easy to review and undo.

## Why one watch list for the script and the agent

The script reads its terms from `references/watch-list.md`, the same
file the agent reads for the plain alternatives. One list means the
two can never disagree about what a term is or how to fix it.

## Why templ files are supported

Page text is often written in templates, not in Markdown. templ is
a widely used way to write HTML pages in Go. Its markup is regular
enough for a script to find the text between tags. Checking only
that text keeps Go code, attributes and expressions out of the
results. Other template languages can be added the same way.
