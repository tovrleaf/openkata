# Docker Image Vulnerability Scanner Skill

## Problem Description

A platform security team is packaging a vulnerability scanning workflow into a reusable agent skill. The workflow helps engineers scan Docker images for known CVEs, triage results by severity, and produce a remediation report. Engineers currently run this manually using a mix of CLI tools and shell scripts; the goal is to capture it as a skill so any agent can execute it consistently when asked.

Create a skill for this Docker image scanning workflow. The skill should guide an agent through: identifying the target image, running a scan, interpreting results by severity tier, and producing a structured remediation summary.

## Output Specification

Produce the following files:

1. `vuln-scan-skill/SKILL.md` — the complete skill definition
2. `validation-report.md` — a document listing the prompts you used to test the skill, with one section per prompt containing:
   - The prompt text
   - Whether the skill should trigger (yes/no) and why
   - Which part of the skill body would guide the agent's first action
3. Any additional skill support files needed (place in appropriately named subdirectories under `vuln-scan-skill/`)
