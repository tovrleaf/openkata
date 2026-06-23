---
name: go-testing
description: >
  Go testing conventions applied to all _test.go files.
  Standard library only, table-driven, behavior-focused.
metadata:
  version: "1.0.0"
  tags: "category:conventions, language:go"
---

# Go Testing

Write tests that are focused, readable, and use only the
standard library.

## Framework

- Standard library only — no testify or external test
  frameworks.
- Run tests: `go test ./...`
- Run with race detection: `go test -race ./...`

## Structure

- Test files live next to source: `handlers.go` →
  `handlers_test.go`
- Test fixtures go in `testdata/` directories.
- HTTP handlers tested with `net/http/httptest`.

## Table-Driven Tests

- Use table-driven tests with named subtests:

  ```go
  tests := []struct {
      name string
      input string
      want  string
  }{
      {"empty input", "", ""},
      {"single word", "hello", "HELLO"},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          got := Func(tt.input)
          if got != tt.want {
              t.Errorf("Func(%q) = %q, want %q", tt.input, got, tt.want)
          }
      })
  }
  ```

## Failure Messages

- Include function name, inputs, got, and want:
  `t.Errorf("Func(%v) = %v, want %v", input, got, want)`
- Never use generic messages like "unexpected result".

## Error Assertions

- Use `errors.Is()` and `errors.As()` — never string
  comparison on error messages.

## Test Helpers

- Test helpers must call `t.Helper()` as the first
  statement so failures point to the caller.

## Scope

- Test behavior, not implementation details.
- Omit tests for trivial logic (simple getters, constant
  returns). Prioritize business logic and conditional
  branches.

## Coverage

- Do not aim for 100% coverage. Cover decisions and
  branches that matter.
