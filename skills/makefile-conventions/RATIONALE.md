# Rationale

## Why Make as interface, scripts as implementation

Makefiles become unreadable when they contain logic.
Delegating to scripts keeps make targets as a
discoverable command menu while scripts handle
complexity. Each can be tested independently.

## Why ban the catch-all target

`%: @:` silently swallows typos. `make tset` succeeds
with no output instead of failing. Explicit targets
fail loud on typos, which is correct behavior.

## Why bare domain shows help

When `make chat` has multiple subcommands, running it
bare should show what's available — not silently do
nothing or pick a default. Discoverability over
convenience.

## Why modular mk/ includes

A single Makefile grows unwieldy past 100 lines.
Splitting by domain (chat.mk, dev.mk, skills.mk)
lets teams own their section without merge conflicts.
The indirection cost is low; the organization value
is high.
