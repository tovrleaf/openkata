# Notification System Refactor Plan

## Background

The current `NotificationService` is a monolith that handles email, SMS, and push
notifications in a single class. It is synchronous and tightly coupled to the
order service. As volume grows, we need better isolation and scalability.

## Proposed Design

Split the notification system into three bounded contexts:
- **Email** (`src/notifications/email/`) — transactional emails via SendGrid
- **SMS** (`src/notifications/sms/`) — time-sensitive alerts via Twilio
- **Push** (`src/notifications/push/`) — mobile push via Firebase

Move from synchronous calls to an event-driven model. The order service will
emit events; each notification context will consume relevant events independently.

## Technology Choice

Use Redis Streams as the internal message queue. We considered RabbitMQ and an
in-process EventEmitter but chose Redis because we already use Redis for caching
and want to avoid a new infrastructure dependency.

## Open Questions

- Should each notification context have its own retry logic, or share a common one?
- Who owns the "notification preferences" data — the notification service or the user service?
- Should SMS and push share a single context since they're both mobile-channel?
