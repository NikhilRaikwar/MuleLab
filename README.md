# MuleLab

> **MuleLab asks a simple question: should an AI agent keep its job?**

Give the sandbox a commerce objective. Specialist agents propose measurable experiments, operate through schema-validated tools, stay inside deterministic budgets and permissions, and are evaluated on business outcomes. Agents that fail their stop conditions are retired. Every model call, tool call, approval, metric, retry, and retirement decision is traceable.

MuleLab is a Sticker Mule-inspired autonomous commerce-agent sandbox built only from public product/workflow information. It is not affiliated with Sticker Mule and uses no private/internal Sticker Mule data or APIs.

## Demo

- Live demo: _pending deployment_
- Recruiter route: `/demo`
- Screenshots: _added after the UI is verified_

The curated flow demonstrates meaningful autonomous work, explicit tools, bounded authority, measured outcomes, deterministic retirement, model evaluation, and a complete evidence trace in under two minutes.

## Architecture

Next.js/TypeScript presents the public experience. A Go/gqlgen API runs a bounded agent workflow, OpenRouter gateway, typed tool registry, deterministic policy and simulation engines, eval runner, and trace recorder. PostgreSQL persists runs, proposals, approvals, tool/model calls, outcomes, and evals. API and web target scale-to-zero GCP Cloud Run; Neon is the zero-cost-first Postgres target.

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Agent lifecycle

`PROPOSED → POLICY_CHECKED → AWAITING_APPROVAL/APPROVED → RUNNING → EVALUATING → WINNING/LOSING/PAUSED/RETIRED/FAILED`

The LLM proposes hypotheses and typed actions. Deterministic code owns authorization, approvals, execution, metrics, thresholds, idempotency, and retirement.

## Tools and safety

Every agent action passes through a typed allowlisted tool with permission checks, exact argument validation, idempotency, tracing, and bounded failure handling. Business/customer/model content is untrusted and cannot expand authority. See [docs/SECURITY.md](docs/SECURITY.md).

## Evals and routing

Golden cases exercise tool correctness, budgets, lifecycle rules, idempotency, malformed output, failures, and prompt injection. Deterministic truth is graded deterministically; optional LLM judges handle only subjective rationale quality. `openrouter/free` is the default live route, with no silent paid fallback and a clearly labeled deterministic demo fallback. See [docs/EVALS.md](docs/EVALS.md).

## Local setup

See [docs/SETUP.md](docs/SETUP.md) and copy `.env.example` without committing the resulting `.env`. Live AI needs an OpenRouter key; the deterministic demo does not.

## Deployment

See [docs/DEPLOY_GCP.md](docs/DEPLOY_GCP.md). No cloud resource is provisioned without an explicit cost review and approval.

## Deliberate limitations

- Fictional seeded business data; no real Sticker Mule integration.
- No authentication or claim of production-grade multi-tenancy in v0.1.
- No real email/SMS, money movement, ad spend, or arbitrary code/URL execution.
- Free-model availability is not an SLA; the recruiter flow remains usable through a labeled deterministic adapter.
- Model arena results appear only when actually executed; benchmarks are never fabricated.

## Project documentation

- [Setup audit](docs/SETUP_AUDIT.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Setup](docs/SETUP.md)
- [Evals](docs/EVALS.md)
- [Security](docs/SECURITY.md)
- [Decisions](docs/DECISIONS.md)
- [GCP deployment](docs/DEPLOY_GCP.md)
