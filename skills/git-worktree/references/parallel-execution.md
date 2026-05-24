# Parallel Execution Pattern

Run batch operations across multiple tasks using one
worktree per task.

## Example: Eval All Skills

```bash
# Create one worktree per task
git worktree add .worktrees/eval-skill-a -b eval/skill-a
git worktree add .worktrees/eval-skill-b -b eval/skill-b
git worktree add .worktrees/eval-skill-c -b eval/skill-c

# Run tasks in parallel (each in its own directory)
cd .worktrees/eval-skill-a && tessl eval run skills/skill-a/ &
cd .worktrees/eval-skill-b && tessl eval run skills/skill-b/ &
cd .worktrees/eval-skill-c && tessl eval run skills/skill-c/ &
wait

# Merge results back
git merge eval/skill-a eval/skill-b eval/skill-c

# Cleanup
git worktree remove .worktrees/eval-skill-a
git worktree remove .worktrees/eval-skill-b
git worktree remove .worktrees/eval-skill-c
git branch -d eval/skill-a eval/skill-b eval/skill-c
```

## Multi-Agent Pattern

When delegating to multiple agents:

1. Create a worktree per agent before spawning
2. Each agent operates in its own `.worktrees/<task>/`
3. Agents commit independently on their branches
4. Coordinator merges results from the main checkout
