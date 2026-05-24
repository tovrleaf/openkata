# TypeScript Interfaces

Enforces consistent naming and declaration of TypeScript interfaces.

## Naming

- Prefix all interface names with `I` (e.g., `IUser`, `IProduct`).
- Use PascalCase for the portion after the prefix (e.g., `IOrderItem`, not `Iorderitem`).
- Do not use the `I` prefix for type aliases — only for `interface` declarations.

## Declaration Style

- Use `interface` for object shapes that may be extended or implemented.
- Use `type` for unions, intersections, and primitives.
- Never use `interface` and `type` interchangeably for the same shape.

## Placement

- Place shared interfaces in `src/types/` — not inline in component or service files.
- Co-locate component-specific interfaces in the same file as the component.
