# Frontend Agent Profile

## Problem Description

A mid-sized product team has a TypeScript web application split between a React frontend (`src/components/`) and a Node.js API layer (`src/api/`). They already have an agent that owns the API layer. Now they want a counterpart agent to handle all the React component work — building and refactoring UI components, managing styles, and keeping the frontend consistent.

The team has shared coding standards in the `rules/` directory that all agents are expected to follow, and an existing profile in `profiles/` for the backend agent. Before creating the new profile, explore the project structure to understand what directories exist and what rules are already in place — this will help you scope the new agent correctly and reference the right constraints without duplicating information that's already captured elsewhere.

The project files are in the current directory.

## Output Specification

Create an agent profile for the frontend role. Save it as a markdown file in the appropriate profiles location for this project. Name it after the agent role using a short, lowercase name.
