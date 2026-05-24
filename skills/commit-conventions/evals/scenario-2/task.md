# Committing a Breaking API Change

## Problem/Feature Description

A platform team is retiring the `v1/charge` endpoint of their payment service. All API consumers have been migrated to the new `v2/payment` endpoint over the past quarter, and the old endpoint is now being removed entirely. This is a breaking change — any client that hasn't migrated will stop working when this is deployed.

The removal has been implemented in `inputs/src/routes/payments.js`, which has already been updated to strip out the old endpoint handler and its associated middleware. There is also a project root that may contain commit configuration files — you should check for any existing commit conventions before deciding how to format the message.

The team expects the commit history to clearly communicate the nature and impact of this change so that downstream consumers reading the changelog understand exactly what broke and what they need to do. The commit must be production-ready — no placeholder or work-in-progress markers.

## Output Specification

Initialize a git repository using the files in `inputs/` as the starting state, then commit the breaking change in `inputs/src/routes/payments.js` according to best practices. Save the output of `git log -1` to `commit-log.txt` and the full output of `git log -1 --format="%B"` (just the commit message body) to `commit-message.txt`.
