# Payment Service Integration Plan

## Overview

We need to integrate Stripe as our payment provider. The checkout flow will
collect card details, the system will handle payment and then the notification
service will send a confirmation to the customer.

## Domain Concepts

- **Checkout**: The process where the customer enters payment details
- **Payment**: When money moves from customer to us
- **Transaction**: A record of the payment attempt
- **Order**: The customer's purchase request

## Flow

1. Customer completes checkout
2. System initiates payment via Stripe API
3. Stripe confirms payment
4. System creates a transaction record
5. Notification service sends confirmation email

## Open Questions

- Should we store card details locally or rely entirely on Stripe tokens?
- What happens if the notification fails after payment succeeds?
- Who owns the transaction record — the payment service or the order service?
