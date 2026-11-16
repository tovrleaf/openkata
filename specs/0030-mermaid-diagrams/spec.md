---
status: Draft
depth: Standard
---

# Diagrams from Mermaid Source

## Story

As a docs contributor, I want to write diagrams as Mermaid text in
the Markdown page. Then adding or changing a diagram is a text
edit, and the diagram shows in any Markdown viewer.

## Background

- Spec 0025 draws diagrams by hand as SVG files that the docs
  renderer inlines into the page. That works without JavaScript
  and follows the site theme, but every change means moving
  coordinates
- More diagrams are expected in the docs
- A `mermaid` code block shows as a diagram on GitHub and in most
  editors, so the raw Markdown reads well locally. An image link
  to an SVG file shows only the picture, and the file is not
  editable as text
- Mermaid v13 renders SVG in Node with no DOM or browser. Its
  labels are SVG text, and custom styling goes through
  `--mermaid-*` CSS variables. It was in alpha on 2026-11-16
  (13.0.0-alpha.3), marked "Don't use it in production". See
  <https://mermaid.js.org/intro/v13-alpha.html>
- Earlier Mermaid versions need Chromium to render outside a
  browser, which spec 0025 rejected

## Start condition

Work starts when Mermaid v13 has a stable release, or when the
Mermaid team brings browser-free rendering to v12.

## Requirements

- Diagrams are `mermaid` code blocks in the docs Markdown
- Pages ship SVG, not Mermaid's JavaScript: no Mermaid code runs
  in the visitor's browser
- Rendered diagrams follow the site theme in both themes
- Deploys and CI do not need Node to build the site. Rendering
  happens in a `make` target, and its output is committed
- The strict docs loader fails when a `mermaid` block has no
  rendered SVG, or when the SVG is older than the block
- The hand-drawn SVGs from spec 0025 are replaced with `mermaid`
  blocks, unless a diagram reads better hand-drawn

## Out of Scope

- Rendering Mermaid in the browser
- Diagram types other than flowchart, unless the docs need them

## Open Questions

- Can the `--mermaid-*` variables map to the site's design
  tokens, or does the output need rewriting?
- Where do rendered SVGs live? One idea: named by a hash of the
  block's source, next to the page, so a changed block is easy
  to detect
- Do hand-drawn SVGs stay supported next to Mermaid blocks?
