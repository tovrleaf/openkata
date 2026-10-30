---
title: How we use evals
description: How skill evals are written and run with the local runner and its Docker sandbox, how they are scored, and where the results go.
order: 4
---

# How we use evals

An eval checks that an agent following a skill really does what the
skill asks. We write a few scenarios for each skill, run them with a
local runner called `openkata-eval`, and keep the score. That score
is the effectiveness shown on the skill's Benchmarks tab. Only skills
in `skills/` have evals.

## What a scenario is

A scenario is a folder inside the skill, at
`skills/<name>/evals/scenario-<something>/`. The runner finds a
scenario only if its folder name starts with `scenario-`.

| File | What it holds |
|------|---------------|
| `task.md` | The task given to the agent. Required. |
| `criteria.json` | The checklist the judge scores. Required. |
| `scenario.json` | Optional settings, described below. |
| `inputs/` | Optional files the task needs, such as a diff or a config. |

### The checklist

`criteria.json` has a `context` line saying what the scenario tests
and a `checklist` of criteria. Each criterion has a `name`, a
`description` that says what to look for, and a `max_score`.

```json
{
  "context": "What the scenario tests",
  "type": "weighted_checklist",
  "checklist": [
    {
      "name": "imperative-mood",
      "description": "The commit header uses the imperative mood",
      "max_score": 12
    }
  ]
}
```

The runner reads `context` and `checklist` and ignores `type`. By
convention the `max_score` values of a scenario add up to 100, and
`scripts/check-eval-weights.sh` checks that. The runner itself works
with whatever weights it finds.

Each criterion is judged as pass or fail, so its description has to
be something observable. "Header uses imperative mood" can be
judged. "Well-written commit message" cannot. The task must not name
the skill or hint at what it expects.

### Scenario settings

`scenario.json` can set three things:

| Key | Meaning | Default |
|-----|---------|---------|
| `description` | A short label shown in results | empty |
| `sandbox` | Run the agent in a Docker container | `true` |
| `network` | The Docker network for this scenario | the runner's setting |

A scenario without a `scenario.json` is sandboxed and has no
description. Skills that only produce text can set `"sandbox": false`
and then need no Docker.

## Running the evals

```bash
make eval-image                              # once, for sandboxed scenarios
make eval-local skills/<name>                # all scenarios of a skill
make eval-local skills/<name>/evals/scenario-<x>   # one scenario
```

`make eval-local` builds `bin/openkata-eval` if it is missing, then
runs it with the path you give. Running a whole skill is the normal
case. Running one scenario is for debugging a failure, and it always
exits with status 0.

The runner needs `SKILL.md` in the skill folder. It also needs
`kiro-cli` on the `PATH` when the backend is `kiro`, which is the
default, and Docker running when any scenario is sandboxed.

### Settings

Settings are resolved in this order, and a later one wins: built-in
defaults, then a `.openkata-eval.yaml` file in the working directory,
then environment variables, then command-line flags. The repository
does not include a `.openkata-eval.yaml`, so the defaults apply
unless you add one.

| Flag | Meaning | Default |
|------|---------|---------|
| `--backend` | `kiro` or `http` | `kiro` |
| `--model` | Model for the agent | `claude-sonnet-4.6` |
| `--judge-model` | Model for the judge | the agent model |
| `--threshold` | Overall pass percentage | `95` |
| `--output` | File or folder for the JSON report | see Results |
| `--model-label` | Readable model name stored in the report | none |

The file sets more than the flags do: `concurrency` (how many
scenarios run at once, default 2), `direct_timeout` (120 seconds),
the HTTP backend's `http.base_url` and `http.api_key`, and the Docker
`sandbox.image` (default `openkata-eval:latest`), `sandbox.timeout`
(300 seconds) and `sandbox.network` (default `host`). The
environment variables cover the backend, models, threshold, HTTP
settings and direct timeout: `OPENKATA_EVAL_BACKEND`, `_MODEL`,
`_JUDGE_MODEL`, `_THRESHOLD`, `_HTTP_BASE_URL`, `_HTTP_API_KEY` and
`_DIRECT_TIMEOUT`. Concurrency and the Docker settings come only from
the file.

The `http` backend calls an OpenAI-compatible
`<base_url>/chat/completions` endpoint, with the key as a bearer
token.

## What happens in a run

For each scenario, with up to `concurrency` running at once:

1. **Build the prompt.** The skill's `SKILL.md`, then its
   `references/` (without `ACKNOWLEDGMENTS.md`), `scripts/` and
   `assets/`, then `task.md`, then each file in `inputs/` as a code
   block.
2. **Run the agent.** Either directly or in a sandbox, as described
   below.
3. **Judge the result.** A second model checks the result against
   the checklist.
4. **Score it.**

### Direct scenarios

A scenario with `"sandbox": false` sends the prompt straight to the
model, with no tools. With the `kiro` backend the runner starts
`kiro-cli chat` with the chosen `--model`, in a temporary folder. If
the model is throttled, the call is tried up to three times, waiting
1, 2 and then 4 seconds.

### Sandboxed scenarios

A sandboxed scenario runs in a fresh Docker container made from
`Dockerfile.eval`, which is Ubuntu 24.04 with `git`, `curl`, `jq` and
`kiro-cli`. For each scenario the runner:

1. creates and starts a container on the scenario's network (the
   `host` network unless set otherwise), with your `~/.aws` folder
   mounted read-only and `/workspace` as the working folder
2. copies the skill (`SKILL.md`, `references/`, `scripts/` and
   `assets/`) into `/workspace/skill/`, and `task.md` and the
   contents of `inputs/` into `/workspace/`
3. writes an `eval-agent` configuration that allows the file, bash,
   grep and glob tools and exposes the skill files to the agent
4. takes a snapshot of the workspace
5. runs `kiro-cli chat --agent eval-agent --trust-all-tools
   --no-interactive` with the prompt, limited to the sandbox timeout
6. takes a second snapshot, and removes the container

The judge is given the agent's output and the difference between the
two snapshots: the contents of the added and modified files, and the
names of the deleted ones. That is how a scenario can check what the
agent did to files, not only what it said.

Two things follow from how the sandbox is started. It always uses
`kiro-cli`, whatever `--backend` says. And it does not pass `--model`,
so the agent runs under `kiro-cli`'s own default model. The `--model`
setting applies to the judge and to scenarios with `"sandbox": false`.

### The judge

The judge gets the agent's result, the scenario's `context` and every
criterion with its name, points and description. It must answer with
a JSON array with one verdict for each criterion: its `name`,
`pass` and a short `reason`. If the answer cannot be parsed, the
runner asks once more. A result is rejected as suspicious when the
agent's response is under 20 characters and every criterion passed.
Setting `OPENKATA_EVAL_DEBUG` prints the raw answer when parsing
fails.

## Scoring

- A criterion is all or nothing. A pass earns its full `max_score`
  and a fail earns 0. A criterion the judge did not mention counts
  as a fail.
- A scenario's percentage is its points divided by its maximum. A
  scenario passes only at 100%.
- The overall result is the average of the scenario percentages. It
  passes at or above the threshold, 95% unless you change it.
- A scenario that cannot be completed, for example because the
  judge fails, is scored 0 and shows the error in its description.

Running a whole skill exits with status 1 when the overall result is
below the threshold. The runner prints each scenario's criteria with
their points, then the overall result.

## Results

Running a whole skill writes a JSON report to
`skills/<name>/evals/results/<model>/<time>.json`, where `<model>` is
the `--model` value and the time is in UTC, such as
`2026-10-27T120000`. With `--output`, the report goes to that file or
folder instead, with a timestamp added to the name. A single
scenario writes nothing unless you give `--output`.

The report holds the skill name, its `skill_version`, the time, the
backend and models, the threshold, every scenario with its criteria
and points, the overall percentage and whether it passed. The
`skill_version` is the first released version heading, such as
`## [1.0.1]`, in the skill's `CHANGELOG.md`.

### From results to the site

`make versions` reads these reports when it builds `versions.json`
from the repository, as described in
[where artifacts are stored](storage.md). For each model it takes
the newest report and keeps it only if its `skill_version` is the
skill's current version, so a report for an older version never
shows up. The overall result becomes that model's effectiveness, and
`--model-label` becomes its name on the site. These are the figures
on the Benchmarks tab, described in
[using the catalog](../guides/using-the-catalog.md).

The `evals/` folder is not part of an installed skill.

## How we work with them

The project's eval workflow for a skill:

1. Write the scenarios, and check the weights add up.
2. Run the whole skill and read which criteria failed and why.
3. Change `SKILL.md` to fix the behavior, and run again, until the
   overall result is 95% or more.
4. Run `tessl skill review --optimize` for a final polish, and run
   the evals again to confirm nothing regressed.

The local runner is what we use for scoring. We do not use
`tessl eval run` for it, because a mismatch in the model makes its
results unreliable.
