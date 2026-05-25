# Documenting the Async Email Delivery Decision

## Problem/Feature Description

NotifyCore is a SaaS billing and subscription management platform. Its API sends transactional emails — order confirmations, payment receipts, and password resets — synchronously inside the request handler using a direct SMTP call. At low volume this was acceptable, but as the customer base has grown, slow email delivery is delaying API responses and inflating timeout rates. Under load, a single unresponsive mail server can cause an entire request to hang.

The infrastructure team has approved a move to queue-based asynchronous delivery for the platform's own transactional email flows. The change is deliberately narrow: the platform's own email service handles only order confirmations, payment receipts, and password resets. Marketing campaigns, weekly digest emails, and promotional blasts are owned by a separate marketing engineering team and run through a third-party email marketing vendor; those are not part of this project. SMS and push notifications are likewise out of scope.

The project files are in `inputs/`. There are no existing architectural decision records in this project. Create the ADR documenting the decision and save it in the standard location within `inputs/`.

## Output Specification

- Create a new ADR documenting the async email delivery decision, placed in the correct location within `inputs/` with a properly formatted filename.
- Fill all sections with substantive content — no placeholder text remaining.
- Write a `delivery-decision-summary.md` at the root of your working directory containing: the ADR filename, the chosen async delivery approach, and one sentence describing what was explicitly excluded from this decision's scope.
