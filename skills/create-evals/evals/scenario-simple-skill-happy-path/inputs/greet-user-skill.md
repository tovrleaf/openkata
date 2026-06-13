---
name: greet-user
description: >
  Generates personalized greeting messages based on user
  locale and context. Use when the user says "greet",
  "welcome message", or needs a locale-aware salutation.
metadata:
  version: "1.0.0"
  tags: "category:communication"
---

# Greet User

Generate locale-aware greeting messages.

## Workflow

1. **Detect locale** — Check the user's locale from
   environment or explicit input. Default to `en-US`.

2. **Select template** — Choose the appropriate greeting
   template for the locale. Fall back to English if the
   locale is unsupported.

3. **Format output** — Substitute the user's name and
   context into the template. Output the final greeting.

## Conventions

- Always include the user's name if provided
- Use formal register for business contexts
- Use informal register for casual contexts
- Never use emoji in formal greetings

## Boundaries

- DOES generate greeting text
- Does NOT send messages
- Does NOT modify user preferences

## Common Failures

- **Wrong register** — using casual tone in a business
  email context.
- **Missing fallback** — crashing on unsupported locales
  instead of falling back to English.
