# Analytics Platform Deployment Architecture

**Version:** 0.4  
**Author:** Platform Engineering  
**Status:** Draft — pre-ratification review

---

## Overview

This document describes the deployment architecture for Prism, a new multi-tenant real-time analytics platform. Prism ingests event streams from customer applications, processes them, and surfaces dashboards with sub-second latency.

---

## 1. Cloud Provider

We will deploy on Google Cloud Platform (GCP). The GCP Console, IAM, and billing tooling are already used for our data warehouse. Expanding to GCP for compute avoids cross-cloud complexity.

**Decision:** Use GCP.

---

## 2. Container Orchestration

We will use Google Kubernetes Engine (GKE) to manage all service workloads. Teams already have internal Kubernetes expertise from the legacy platform. GKE supports autoscaling and rolling deployments out of the box.

**Decision:** Use GKE for container orchestration.

---

## 3. Primary Database

All tenant configuration, user data, and aggregated metrics will be stored in Cloud SQL (PostgreSQL). PostgreSQL is the team's preferred relational database and Cloud SQL provides managed backups, failover, and patching.

**Decision:** Use Cloud SQL (PostgreSQL) for primary storage.

---

## 4. Event Stream Processing

Incoming events will be ingested via Cloud Pub/Sub and processed by a fleet of stateless workers. Workers decode, validate, and fan out events to downstream consumers.

**Decision:** Use Cloud Pub/Sub for event ingestion.

---

## 5. In-Memory Cache

A Redis cluster will be used to cache aggregated dashboard metrics and per-tenant rate-limit counters. The existing ops team is comfortable with Redis.

**Decision:** Use Redis for caching.

---

## 6. Deployment Regions

The initial launch will target a single region: `us-central1`. A second region (`eu-west1`) will be added in Q3 based on customer demand from Europe.

**Decision:** Single-region launch; expand to multi-region in Q3.

---

## 7. Service Communication

Internal services will communicate over gRPC. gRPC provides strong typing via protobuf definitions, which simplifies contract management between teams.

**Decision:** Use gRPC for internal service communication.

---

## Open Questions

- Multi-tenancy isolation model (shared schema vs per-tenant schema) is TBD.
- SLA targets for dashboard load time have not been finalized.
- Disaster recovery runbook is not yet written.
