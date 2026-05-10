# Design System

Always-on conventions for token-based CSS design systems.
Apply when writing or modifying CSS, templates, or page
layouts.

## Token Usage

- Use `var(--token)` for all values. Never hard-code
  colors, spacing, or font stacks.
- Semantic color tokens (`--color-ink`, `--color-surface`,
  `--color-border`, `--color-accent`) for themed elements.
- Raw palette tokens only for fixed-color elements like
  badges.

## Themes

Each theme defines semantic color tokens via
`[data-theme]` selectors. New components must work across
all themes — never reference theme-specific colors
directly.

## File Structure

- `tokens.css` — structural tokens (spacing, typography,
  layout). No colors.
- `themes/*.css` — one file per theme, semantic colors
  only.
- `base.css` — reset, body, links.
- `typography.css` — headings, code, decorative elements.
- `components/*.css` — one file per component.
- `layouts/*.css` — grid and layout patterns.
- `style.css` — imports all files, page-specific styles.

New components go in `components/`. New layouts go in
`layouts/`. Add the `@import` to `style.css`.

## Borders

- Use solid borders. No dashed or dotted.
- Use border color tokens for structural lines.
- Use accent color tokens for emphasis borders.
- No inline styles. All styling lives in CSS files.

## Naming

- BEM-style for components: `.card`, `.card__title`.
- Utility classes for decorative patterns: `.comment`,
  `.cursor`, `.badge`.
- Page-specific classes for unique elements: `.hero`,
  `.nav-brand`, `.footer`.

## Fonts

- Define font stacks as tokens (`--font-body`,
  `--font-mono`).
- Use body font for prose, mono for headings and
  technical elements.
- Vendor fonts to the project rather than relying on
  external CDNs in production.

## Layout

- Constrain content width via token (`--max-width`).
- Full-bleed elements break out of max-width using
  viewport width techniques.
