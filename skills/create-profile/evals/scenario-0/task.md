# Documentation Agent Profile

## Problem Description

A growing engineering team has built a TypeScript web application with a `docs/` directory for architecture guides, API references, and onboarding materials. The team wants to introduce a dedicated AI agent responsible for maintaining and expanding this documentation — keeping it accurate, well-structured, and up to date with the codebase.

The project is under active development by multiple contributors, so it's important that the documentation agent stays strictly within its lane: it should handle the `docs/` folder and nothing else. When documentation work touches code (e.g., updating an API example requires a code change), the agent should flag it and hand off rather than attempting the change itself.

The project files are in the current directory. Explore the structure to understand what directories exist and what scope makes sense for the agent before creating the profile.

## Output Specification

Create a sensei agent profile for the documentation role. The profile should be saved as a markdown file in the appropriate location for profiles in this project. Name the profile file after the agent role.
