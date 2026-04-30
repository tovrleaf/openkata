---
name: create-skill
description: >
  Creates new agent skills following the Agent Skills specification.
  Investigates the repo for conventions, designs the skill around
  progressive disclosure, writes SKILL.md with effective trigger
  descriptions, and validates with representative prompts. Use when
  the user wants to create a skill, build a SKILL.md, turn a
  workflow into a reusable skill, teach the agent a new task,
  scaffold a new agent capability, has a repeated workflow they
  want to codify, is frustrated by inconsistent agent behavior, or
  wants to package expertise for a team.
---

# Create Skill

Create agent skills that are portable, easy to trigger, and cheap
to load. A skill is a folder containing a SKILL.md file with YAML
frontmatter and markdown instructions.

## Workflow

1. **Intake (mandatory gate)** — Understand what the skill
   should do. Ask at least 3 targeted questions before
   drafting anything. Collect:
   - A short name (lowercase, hyphenated)
   - What the skill enables the agent to do
   - When it should activate (trigger conditions)
   - What success looks like

   Summarize your understanding and get explicit confirmation
   before proceeding. Do not write SKILL.md until the user
   confirms.

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

4. **Design the package** — Use progressive disclosure. For
   iterative development, draft in a temp directory (e.g.,
   `/tmp/skills/skill-name/`) to avoid cluttering the repo
   until the skill is validated.

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
     say. Include both *actions* (what the user asks to do) and
     *situations* (what the user is experiencing). Keep workflow
     details out of the description.

     Bad: `Follows a 7-step process to generate SKILL.md files
     with YAML frontmatter containing name and description
     fields.`

     Good: `Creates agent skills. Use when the user wants to
     build a SKILL.md, turn a workflow into a reusable skill,
     has a repeated workflow they want to codify, or is
     frustrated by inconsistent agent behavior.`
   - **Body is procedural and imperative.** Tell the agent
     exactly how to proceed. Don't restate trigger criteria
     from the description — if the body has a "When to use"
     section, it's duplicating the description.
   - **Explain the why.** Reasoning is more effective than rigid
     `ALWAYS`/`NEVER` rules.
   - **Be concise.** Remove explanations of concepts the model
     already knows. Trim examples that restate the rules.

   See [example-skill.md](references/example-skill.md) for a
   complete finished skill demonstrating these principles.

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

## Example Scenario

User: "I keep doing the same database migration steps manually,
can we turn that into a skill?"

1. Skill investigates repo — finds Flyway config, existing
   migration scripts, deploy docs
2. Asks: "Should this skill also handle rollback, or just
   forward migrations?"
3. Creates `migrate-database/SKILL.md` with detect → plan →
   execute → verify workflow
4. Moves Flyway-specific edge cases into `references/`
5. Validates with positive prompt ("run the migration") and
   negative ("change the database schema" — different job)

## Common Failures

- **Description too vague to trigger** — "helps with code" won't
  activate. Include specific actions and situations.
- **Description leaks workflow** — the agent reads the summary
  and skips the body, following a shortcut instead of the full
  procedure.
- **Body too abstract to act on** — "investigate the problem"
  isn't actionable. "Run `git log --oneline -20` to check
  recent patterns" is.
- **No investigation step** — the skill asks the user questions
  it could have answered by reading the repo.

## Quality Checklist

Use the
[skill design checklist](references/skill-design-checklist.md)
before finalizing. For detailed prompt-testing guidance, see
[skill validation](references/skill-validation.md). For trimming
context cost, see
[token optimization](references/token-optimization.md).
