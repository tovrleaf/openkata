---
title: Using the catalog
description: Browse, search and read the skills, rules and profiles on openkata.dev, and download them.
order: 4
---

# Using the catalog

The site lists every skill, rule and profile. You can browse each
type, search all of them at once, read an artifact before you use
it, and download it.

## Browse by type

The [skills](/skills/), [rules](/rules/) and [profiles](/profiles/)
pages each list their artifacts alphabetically. An artifact that
has no release yet (version `0.0.0`) is not listed.

Each entry shows:

- its name and latest version
- its download count, for example `42 ↓`
- for skills with eval results, the best effectiveness score and
  the model that earned it
- its tags, each a link that searches the catalog for that tag

Click an entry to expand it and read its description. An expanded
skill also shows the score for each model. **Open** takes you to
the artifact's own page.

## Search

The [catalog](/catalog/) searches skills, rules and profiles
together. Type in the box and the results update as you type. The
address bar keeps your query (`/catalog/?q=...`), so you can share
a link to a search.

The magnifying glass in the top bar opens the same search in an
overlay, from any page. It shows the first seven matches. Press
Enter to continue on the catalog page, or choose **View all** when
more than seven skills, rules and profiles match, and press Escape
to close it.

The overlay also finds pages of this documentation. They are marked
`docs`, and a page matches on its title, its description and its
headings. Skills, rules and profiles always come first, in order of
how well they match, and documentation pages always come last, in the
same order, however well they match. The catalog page and **View all** list
skills, rules and profiles only. A search with nothing typed lists
artifacts only, and a term with a colon, which is a tag, never
matches a documentation page.

### What a search matches

Search is case-insensitive and looks at an artifact's name, its
tags and its description. Separate terms with spaces, and every
term must match. A match in the name counts most, then one in the
tags, then one in the description, so the best matches come first.
Ties are listed alphabetically. With nothing typed, every artifact
is listed alphabetically.

### Searching by tag

A term that contains a colon is a tag. It must equal one of the
artifact's tags exactly, such as `tool:git`, `language:go` or
`category:conventions`. `category` alone is not a tag, so it
matches by the usual name, tag and description rule instead.

You can mix both kinds of term. `git category:conventions` finds
artifacts that mention git and carry the `category:conventions`
tag. Clicking a tag anywhere on the site runs the same search.

## Read an artifact

Every artifact has its own page, at `/skills/<name>/`,
`/rules/<name>/` or `/profiles/<name>/`. The top of the page shows:

- a version selector. The newest version is marked *Latest*, and
  choosing another one opens that version's page.
- arrows to the previous and next artifact of the same type, in
  alphabetical order
- the download count and a `.tar.gz` link
- the tags

The tabs below hold the details. Some tabs appear only when the
artifact has that content.

| Tab | On | What it shows |
|-----|----|---------------|
| Overview | all types | The description and the artifact's documentation |
| Files | skills and rules | The files of the version you are viewing |
| Rationale | all types, when present | Why the artifact is the way it is |
| Changelog | all types | What changed in each version |
| Benchmarks | skills with scores | Quality score and effectiveness |
| Acknowledgments | skills and rules, when present | Sources the artifact draws on |

The address bar records the open tab, for example `#changelog`, so
a link to a tab opens that tab.

### Benchmarks

The Benchmarks tab shows two kinds of score for a skill:

- **Quality score.** A structural score from tessl's lint of the
  skill: formatting, metadata and the structure of `SKILL.md`. It
  links to the skill's page in the tessl registry.
- **Effectiveness.** How well the skill performs across its eval
  scenarios when an LLM uses it, per model. Each scenario has a
  bar: green at 95% and above, yellow from 70%, and red below that.

## Download

The `.tar.gz` link on an artifact's page downloads the latest
version. For a specific version, use
`/skills/<name>/archive/<version>`, and the same form for rules and
profiles. Extract the archive into the folder your agent reads, as
described in [installing artifacts](installing-artifacts.md).
Each archive download counts toward the download count.
