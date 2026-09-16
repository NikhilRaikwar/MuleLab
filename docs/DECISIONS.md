# Architecture decisions

## ADR-001: One bounded runtime, four agent roles

Use one orchestration runtime with Supervisor, Store, Notify, and Give role profiles. This preserves meaningful specialization while avoiding fake distributed-agent complexity. The runtime owns iteration, time, tool, and model-call budgets.

## ADR-002: Model proposes; code authorizes

Strict JSON Schema is a parsing boundary. Deterministic policy, approvals, tool authorization, simulator metrics, and transactional lifecycle rules are the authority boundary.

## ADR-003: Deterministic demo is always available

The recruiter path cannot depend on free-model availability. A labeled deterministic proposal adapter supplies the curated flow when live AI is unavailable. Live and fixture-backed runs are never presented as equivalent.

## ADR-004: Direct OpenRouter HTTP adapter

Use Go `net/http` behind a provider-neutral interface. This keeps timeouts, retries, response metadata, routing, and paid-model restrictions explicit and avoids coupling the domain runtime to an SDK.

## ADR-005: PostgreSQL with pgx and explicit repositories

Use `pgx/v5`, SQL migrations, and explicit repositories. The domain stays testable with an in-memory repository while the deployed system uses Neon PostgreSQL.

## Local and hosted database parity

The same SQL migration runner and idempotent Cosmic Cats seed execute against PostgreSQL 17 in Docker and are intended for Neon through its direct, non-pooled migration URL. The runtime uses the pooled URL. The local Docker path is verified; hosted Neon creation is currently blocked by the account's Vercel-managed organization policy, not by the application.

## ADR-006: Current stable Go

Target Go 1.27, satisfying the PRD's Go 1.25+ requirement. Current official release documentation lists Go 1.27.1 as the latest patch as of the audit.

## ADR-007: Two Cloud Run services

Deploy web and API separately on Cloud Run, each with minimum instances 0 and maximum 1. This produces a coherent GCP portfolio story and independent scaling at the cost of two cold-start boundaries.

## ADR-008: Product traces plus optional OpenTelemetry

Persist domain-level trace events because recruiters need durable evidence. Emit OpenTelemetry spans optionally; no hosted tracing vendor is mandatory.

## Build note: dependency integrity quarantine

On 2026-09-16, Go's checksum verifier rejected two independently downloaded gqlgen CLI transitive module metadata files. Checks were not disabled and sums were not edited. Core dependency-free packages remain testable, but generated gqlgen files and the full API build are quarantined until a clean module cache/network session can reproduce trusted checksums.
