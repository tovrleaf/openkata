# Rationale

makefile-conventions structures Makefiles as a universal
command interface with modular includes and
self-documenting help.

## Why Make is the interface and scripts are the implementation

Makefiles become unreadable when they contain logic.
Delegating to scripts keeps make targets as a
discoverable command menu while scripts handle
complexity. Each can be tested independently.

## Why the catch-all target pattern is banned

`%: @:` silently swallows typos. `make tset` succeeds
with no output instead of failing. Explicit targets
fail loud on typos, which is correct behavior.

## Why bare domain commands show help

When `make chat` has multiple subcommands, running it
bare should show what's available — not silently do
nothing or pick a default. Discoverability over
convenience.

## Why targets are split into modular mk/ files

A single Makefile grows unwieldy past 100 lines.
Splitting by domain (chat.mk, dev.mk, skills.mk)
lets teams own their section without merge conflicts.
The indirection cost is low; the organization value
is high.
