# Database Decomposition Plan: Monolith to Microservices

## Background

The platform currently runs on a single shared PostgreSQL database serving all application domains (users, orders, products, inventory, payments). As the engineering org grows, each team needs independent deployability and the ability to choose the right data store for their domain.

## Goal

Split the shared database into domain-owned databases, one per microservice. Each service will own its schema and data, with no direct cross-service database queries.

## Approach

We will use the strangler fig pattern over two sprints.

**Sprint 1 — Dual Write**

- Deploy microservices alongside the monolith (each with their own database instance)
- Modify the monolith to write to both the old shared database and the new service databases simultaneously
- Run backfill migration scripts to copy historical data into each new database

**Sprint 2 — Cutover**

- Validate that data in new databases matches the shared database
- Switch read traffic from the monolith's shared database to the microservice databases
- Remove dual-write logic from the monolith
- Deprecate the shared database

## Definition of Done

- All services reading and writing exclusively to their own databases
- No remaining references to the shared database in application code
- Shared database decommissioned

## Team

The platform team will own the migration. Individual service teams will review their own data models.
