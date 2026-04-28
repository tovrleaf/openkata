---
name: create-skill
description: >
  Creates new agent skills following the Agent Skills specification.
  Investigates the repo for conventions, designs the skill around
  progressive disclosure, writes SKILL.md with effective trigger
  descriptions, and validates with representative prompts. Use when
  the user wants to create a skill, build a SKILL.md, turn a
  workflow into a reusable skill, teach the agent a new task, or
  scaffold a new agent capability.
---

# Create Skill

Create agent skills that are portable, easy to trigger, and cheap
to load. A skill is a folder containing a SKILL.md file with YAML
frontmatter and markdown instructions.

## Workflow

1. **Gather intent** — Understand what the skill should do. Collect:
   - A short name (lowercase, hyphenated)
   - What the skill enables the agent to do
   - When it should activate (trigger conditions)
   - What success looks like

2. **Investigate the repo** — Before asking questions, search the
   repo for facts that reduce ambiguity:
   - Existing skills, conventions, and workflow docs
   - Scripts, templates, schemas relevant to the target workflow
   - Tool or dependency requirements
   - Whether the conversation already contains a workflow to
     capture

   Prefer targeted search over broad reading. Don't ask the user
   what you can look up yourself.

3. **Clarify** — Ask only questions that materially affect the
   skill. Push until these are clear:
   - Required workflow steps and their order
   - Required inputs and expected outputs
   - Dependencies on tools, scripts, or services
   - Whether the skill needs `references/`, `scripts/`, or
     `assets/`

4. **Design the package** — Use progressive disclosure:

   ```text
   skill-name/
   ├── SKILL.md        # Metadata + core workflow
   ├── references/     # Detailed docs, loaded on demand
   ├── scripts/        # Executable code
   └── assets/         # Templates, resources
   ```

   - Keep SKILL.md under 500 lines
   - Move bulky detail into `references/`
   - Put deterministic execution in `scripts/`
   - Don't duplicate guidance across files

5. **Write SKILL.md** — Structure:

   ```markdown
   ---
   name: skill-name
   description: >
     What the skill does and produces. Use when the user wants
     to <scenario>, mentions <keyword>, or asks about <topic>.
   ---

   # Skill Name

   One-line stance on what this skill does.

   ## Workflow

   1. **Step** — Description.

   ## Conventions

   - Key constraints.
   ```

   Writing rules:
   - **Description optimizes activation, not teaching.** State
     the job and when to use it in words a user would actually
     say. Include synonyms and paraphrases. Keep workflow details
     out of the description.

     Bad: `Follows a 7-step process to generate SKILL.md files
     with YAML frontmatter containing name and description
     fields.`

     Good: `Creates agent skills. Use when the user wants to
     build a SKILL.md, turn a workflow into a reusable skill,
     or teach the agent a new task.`
   - **Body is procedural and imperative.** Tell the agent
     exactly how to proceed.
   - **Explain the why.** Reasoning is more effective than rigid
     `ALWAYS`/`NEVER` rules.
   - **Be concise.** Remove explanations of concepts the model
     already knows. Trim examples that restate the rules.

   **Complete example** — A finished skill for running the test
   suite looks like this:

   ```markdown
   ---
   name: run-tests
   description: >
     Runs the project's test suite and reports results. Use when
     the user wants to run tests, check if tests pass, verify
     changes don't break anything, or asks about test failures.
   ---

   # Run Tests

   Run the full test suite, surface failures clearly, and suggest
   fixes when the cause is obvious.

   ## Workflow

   1. **Detect the test runner** — Check package.json scripts,
      Makefile targets, or pyproject.toml for the test command.
      Prefer `npm test`, `make test`, or `pytest` in that order.

   2. **Run the suite** — Execute the detected command. Capture
      stdout and stderr.

   3. **Report results** — If all tests pass, confirm with a
      one-line summary. If tests fail, list each failing test
      with its error message and the file:line reference.

   4. **Suggest fixes** — For failures with an obvious cause
      (import error, missing env var, typo), propose a concrete
      fix. For ambiguous failures, ask the user before changing
      anything.

   ## Conventions

   - Never modify test files to make tests pass.
   - Run the full suite unless the user explicitly scopes to a
     subset.
   ```

   This example shows a realistic description with natural trigger
   phrases, a tight four-step workflow, and conventions that
   constrain behavior without being rigid rules.

6. **Validate** — Test the skill with representative prompts:
   - 2–3 realistic positive prompts (things users would say)
   - At least 1 negative prompt (adjacent but shouldn't trigger)
   - Note whether failures come from trigger wording, workflow
     ambiguity, or missing resources

   Skip validation only for trivial skills where the trigger
   surface is obvious.

7. **Acknowledge sources** — If the skill draws on external
   practices, create `references/ACKNOWLEDGMENTS.md` listing
   each source with a link, license, what was adapted, and the
   version it was adopted in.

8. **Confirm** — Show the user the created skill and ask if
   adjustments are needed.

## Quality Checklist

Use the
[skill design checklist](references/skill-design-checklist.md)
before finalizing. For detailed prompt-testing guidance, see
[skill validation](references/skill-validation.md). For trimming
context cost, see
[token optimization](references/token-optimization.md).
