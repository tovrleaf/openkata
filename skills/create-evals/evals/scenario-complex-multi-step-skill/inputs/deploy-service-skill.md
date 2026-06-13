---
name: deploy-service
description: >
  Orchestrates service deployments from test through production.
  Use when the user says "deploy", "release to prod", "push to
  staging", or is preparing a service for release.
metadata:
  version: "2.1.0"
  tags: "category:operations, tool:docker, tool:kubernetes"
---

# Deploy Service

Deploy a service through staging to production with
verification at each gate.

## Workflow

1. **Verify prerequisites** — Check that CI is green,
   dependencies are pinned, and changelog is updated.

2. **Run tests** — Execute the full test suite including
   integration tests. Stop if any fail.

3. **Build artifacts** — Build Docker image, tag with
   git SHA and semver. Push to container registry.

4. **Tag the release** — Create annotated git tag with
   version from CHANGELOG.md.

5. **Push to registry** — Push container image to the
   project's configured registry.

6. **Deploy to staging** — Apply Kubernetes manifests to
   staging namespace. Wait for rollout. Run smoke tests.

7. **Promote to production** — After staging verification,
   apply to production namespace. Monitor for 5 minutes.
   Rollback if error rate exceeds threshold.

## References

- [runbook.md](references/runbook.md) — Detailed
  troubleshooting steps
- [rollback-procedure.md](references/rollback-procedure.md)

## Scripts

- `scripts/build.sh` — Docker build wrapper
- `scripts/deploy.sh` — Kubernetes apply with namespace
- `scripts/smoke-test.sh` — Post-deploy verification

## Boundaries

- DOES build, tag, and deploy containers
- DOES run smoke tests and monitor
- DOES rollback on failure
- Does NOT modify application code
- Does NOT merge PRs
- Does NOT update DNS or load balancer config

## Common Failures

- **Skipping staging** — deploying directly to production
  without staging verification.
- **Stale image tag** — using latest instead of SHA-pinned
  tags causes non-reproducible deploys.
- **No rollback plan** — proceeding without confirming
  rollback works in staging first.
