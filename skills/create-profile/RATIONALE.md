# Rationale

## Why a 40-line limit

Profiles are lenses, not manuals. They scope an agent
to a domain — anything longer means the profile is
trying to teach rather than constrain. Long profiles
also consume tokens every session.

## Why mandatory stop conditions

Without explicit boundaries, agents attempt work
outside their scope silently. A stop condition like
"never modify application code" makes violations
detectable and correctable.

## Why constraints vs design intent

Constraints are non-negotiable ("never push without
confirmation"). Design intent is preferred behavior
("prefer small commits"). The distinction lets agents
prioritize when guidance conflicts.

## Why reference rules instead of repeating them

Repeating rule content in profiles causes drift. When
the rule updates, profiles become stale. References
ensure a single source of truth.
