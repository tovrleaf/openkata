# Commit Style Rule

All commits must follow Conventional Commits format:

```
<type>(<scope>): <short description>
```

Types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`

Examples:
- `feat(auth): add OAuth2 login flow`
- `fix(api): correct null check in user endpoint`
- `docs(readme): update setup instructions`

Scope is optional but recommended when the change is isolated to a module.
Breaking changes must include `BREAKING CHANGE:` in the commit body.
