# PR Automation Script

## Problem/Feature Description

Your team has grown from 3 to 15 engineers over the past year, and the PR creation process has become inconsistent. Some engineers forget to check their branch state before pushing, others accidentally push directly to main, and a few have had their PRs rejected because they didn't run the project's pre-push checks first. The tech lead has asked you to write a reusable shell script that standardizes the entire PR creation process across the team.

The script should handle all the edge cases the team has run into: detecting uncommitted work and prompting about it, refusing to proceed from the wrong branch, checking whether the local branch is out of date with the remote, running any project-specific pre-push checks, and finally pushing and opening the pull request. The script should also handle environments where the GitHub CLI may not be installed.

## Output Specification

Write the script to a file named `create-pr.sh`. It should be a complete, runnable bash script. The script should:

- Be usable from any project directory
- Handle all the branch and state checks before attempting to push
- Create the pull request using the appropriate tooling
- Offer to open the PR in the user's browser after creation

Do not hardcode any repository names or branch names — the script should detect these from the current git context.
