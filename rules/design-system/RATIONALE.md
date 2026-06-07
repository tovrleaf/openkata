# Rationale

## Why all values via var(--token)

Hardcoded values scatter design decisions across
files. Tokens centralize them. Changing a color
means editing one token, not finding every #hex
in the codebase.

## Why tokens.css has no colors

Colors are theme-specific. Structural tokens
(spacing, fonts, radii) are theme-independent.
Separating them means adding a theme = one new
file with color definitions, zero structural changes.

## Why local fonts, never CDNs

External font requests add latency, create privacy
concerns (Google Fonts tracking), and fail offline.
Vendoring fonts removes all three issues at the cost
of a few KB in the repo.

## Why BEM + utility hybrid

Pure BEM is verbose for simple decoration (a color,
a spacing). Pure utility (Tailwind-style) obscures
component structure. Hybrid uses BEM for component
structure and utilities for one-off adjustments.
