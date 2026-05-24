# User Notifications Service — Design Spec

## Overview

The User Notifications Service is responsible for delivering in-app, email, and push notifications to users when triggered by system events across the platform.

## Scope

- **In-scope:** delivering notifications triggered by other services, tracking delivery status, managing user notification preferences
- **Out of scope:** notification template management, notification content creation

## Architecture

Notifications will be stored in a PostgreSQL database. When a new notification event is received from an upstream service, a background worker picks it up from an internal queue and attempts delivery.

For push notifications we will integrate with a third-party provider (to be selected). For email we will use the existing internal email service.

## User Preferences

Users can configure notification preferences per channel (in-app, email, push). If a user has disabled a channel, notifications for that channel will be skipped silently.

## Delivery

The worker will attempt delivery up to 3 times with exponential backoff before marking the notification as `failed`.

## Data Retention

Notifications will be retained in the database for 90 days before being archived.

## Monitoring

Delivery success and failure rates will be tracked per channel. An alert will fire when the failure rate exceeds 10% over a 5-minute window.
