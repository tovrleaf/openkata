# Rollback Procedure

## Automatic Rollback

If error rate exceeds 5% within 5 minutes of deploy:

```bash
kubectl rollout undo deployment/$SERVICE -n production
```

## Manual Rollback

```bash
kubectl set image deployment/$SERVICE \
  app=$REGISTRY/$SERVICE:$PREVIOUS_SHA -n production
```

## Verification After Rollback

1. Check pod status: `kubectl get pods -n production`
2. Run smoke tests: `./scripts/smoke-test.sh production`
3. Verify error rate returned to baseline
