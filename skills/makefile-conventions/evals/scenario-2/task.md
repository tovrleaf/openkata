# Deployment Command Interface for a Multi-Environment Service

## Problem Description

A backend team maintains a containerised API called `harbor` that needs to be deployed to three different environments: `dev`, `staging`, and `production`. Deployments are handled by two scripts — `scripts/docker-build.sh` and `scripts/deploy.sh` — which accept configuration through environment variables: the target environment (`ENV`), the container registry URL (`REGISTRY`), and the version tag (`VERSION`).

Currently, developers run the scripts directly by hand, passing variables on the command line each time. This leads to mistakes: people forget to set `REGISTRY`, use inconsistent `VERSION` strings, or accidentally deploy to the wrong environment. The lead engineer wants a `make`-based command interface that centralises these operations, documents the available commands, and allows the variables to be overridden on the command line when needed.

The build and deploy operations each stand on their own — they do not share subcommands with each other — but together they are the two concerns the interface should expose. Default values for the variables should be sensible (e.g. `ENV=dev`, `VERSION=latest`) so developers can run a quick local build with just `make build`.

## Output Specification

Create a complete Makefile-based command interface for the `harbor` project. The solution should include:

- A root `Makefile` and any supporting files under `mk/`
- Variable definitions with sensible defaults that can be overridden from the command line (e.g. `make deploy ENV=staging VERSION=1.4.2`)

Demonstrate the help system by capturing:
- The output of `make` (no arguments) to `make-help-output.txt`
