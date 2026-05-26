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
metadata:
  version: "1.2.1"
  tags: "category:scaffolding"
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
   repo for:
   - Existing skills, conventions, and workflow docs
   - Scripts, templates, schemas relevant to the target workflow
   - Tool or dependency requirements
   - Whether the conversation already contains a workflow to
     capture

3. **Clarify** — Ask only questions that materially affect the
   skill. Push until these are clear:
   - Required workflow steps and their order
   - Required inputs and expected outputs
   - Dependencies on tools, scripts, or services
   - Whether the skill needs `references/`, `scripts/`, or
     `assets/`

4. **Design the package** — Structure:

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
   ```

   Frontmatter rules:
   - Only use `name` and `description` in frontmatter.
     `version` is not part of the spec and triggers warnings.

   Writing rules:
   - **Description optimizes activation, not teaching.** State
     the job and when to use it in words a user would actually
     say. Include both *actions* and *situations*. Keep workflow
     details out of the description.

     Bad: `Follows a 7-step process to generate SKILL.md files
     with YAML frontmatter.`

     Good: `Creates agent skills. Use when the user wants to
     build a SKILL.md, turn a workflow into a reusable skill,
     or is frustrated by inconsistent agent behavior.`
   - **Body is procedural and imperative.** Tell the agent
     exactly how to proceed. Don't restate trigger criteria
     from the description — a "When to use" section in the
     body duplicates the description.
   - **Use imperative form.** "Do not", "Use", "Run" — not
     "prefer" or "consider".
   - **Be concise.** Terse reminders, not tutorials.
   - **Include a complete example.** One full, copy-paste-ready
     artifact beats scattered snippets.
   - **Include a Boundaries section (mandatory).** List what
     the skill DOES and Does NOT do.
   - **Include a Common Failures section.** List 2–3
     domain-specific mistakes an agent would make without
     this guidance.

   See [example-skill.md](references/example-skill.md) for a
   complete finished skill demonstrating these principles.

6. **Validate** — Test the skill with representative prompts:
   - 2–3 realistic positive prompts (things users would say)
   - At least 1 negative prompt (adjacent but shouldn't trigger)

   Write a brief validation report noting:
   - Which prompts triggered correctly
   - Which failed and why (trigger wording, workflow ambiguity,
     or missing resources)
   - What was fixed based on the failures

   Skip validation only for trivial skills where the trigger
   surface is obvious.

7. **Acknowledge sources** — If the skill draws on external
   practices, create `references/ACKNOWLEDGMENTS.md` listing
   each source with a link, license, what was adapted, and the
   version it was adopted in.

8. **Confirm** — Show the user the created skill and ask if
   adjustments are needed.

## Boundaries

- DOES create skill directories, SKILL.md, references/, scripts/
- DOES validate with representative prompts
- Does NOT publish, tag, or release skills
- Does NOT modify existing skills (use review-skill for that)
- Does NOT create rules or profiles (separate workflows)

## Example Scenario

User: "Turn my database migration steps into a skill."
→ Investigate repo (Flyway config) → ask about rollback scope
→ create `migrate-database/SKILL.md` → validate with prompts.

## Common Failures

- **Description leaks workflow** — the agent reads the summary
  and skips the body, following a shortcut instead of the full
  procedure.
- **Body too abstract to act on** — "investigate the problem"
  isn't actionable. "Run `git log --oneline -20` to check
  recent patterns" is.

## Quality Checklist

Before finalizing, use the
[skill design checklist](references/skill-design-checklist.md),
[skill validation](references/skill-validation.md), and
[token optimization](references/token-optimization.md).
