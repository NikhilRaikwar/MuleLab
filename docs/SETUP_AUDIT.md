# Setup audit

Audited 2026-09-17. The project now has a Go/gqlgen API, Next.js recruiter UI, PostgreSQL migrations, deterministic agent runtime, Docker Compose, and GitHub Actions. This is an implementation status record, not a deployment claim.

| Service | Why needed | Required for MVP? | Connection/plugin needed? | Secret/env required? | Free option? | Status |
| --- | --- | --- | --- | --- | --- | --- |
| GitHub | Source hosting and CI | Required for hosted CI, not local build | Git remote and authenticated `gh`/Git | No app secret | Public repo and Actions allowance | Connected; CI verified green |
| OpenRouter | Live structured agent proposals and arena | Optional for deterministic public demo; required for live AI | OpenRouter account/API key; no Codex plugin | `OPENROUTER_API_KEY` | `openrouter/free`, rate-limited | Live zero-cost inference verified locally |
| GCP | Public Cloud Run deployment | Required for target deployment | `gcloud` login and authorized project | Local auth; project/region env | Cloud Run has free allowance, never guaranteed $0 | Not connected |
| Neon Postgres | Deployed durable state | Required for deployed persistence | Neon account/project and Neon CLI | `DATABASE_URL`, `DATABASE_URL_DIRECT` | Neon Free | Dedicated PostgreSQL 17 project verified |
| Local PostgreSQL | Local integration tests | Recommended | Docker or local Postgres | Local-only URL | Yes | Docker PostgreSQL 17 verified |
| Artifact Registry | Store deployment images | Required for image-based Cloud Run deploy | GCP authentication | No repository credential committed | Small free allowance may apply | Not connected |
| OpenTelemetry collector/vendor | External trace export | No | Optional OTLP endpoint | Endpoint/headers | Local/stdout and several free options | Disabled |
| GitHub Actions | Deterministic CI | Required once hosted | GitHub repository | Optional deployment secrets | Included allowance varies | Not configured yet |

## Local tooling observation

Go 1.27, gqlgen generation, Docker Desktop, Docker Compose, PostgreSQL 17, and the Windows frontend worker temp workaround (`D:\DevTemp`) have been verified. The same migrations and idempotent Cosmic Cats seed have been exercised against local Docker PostgreSQL and hosted Neon PostgreSQL 17. GraphQL has returned Cosmic Cats from both database paths, and the local `postgres`/`api`/`web` Compose recruiter flow has been verified end-to-end.

## Connection actions

No Codex plugin or MCP server is required. Development proceeds with deterministic adapters and fixtures.

1. OpenRouter: the local key enables a bounded optional live summary call. The public demo remains usable without it via a labeled deterministic fallback.
2. Neon: the dedicated MuleLab Free project is configured only through ignored local environment variables; do not put URLs or passwords in documentation.
3. GCP: authenticate `gcloud` and select a billing-enabled project only after the deployment cost review. No cloud resources will be created without explicit approval.
4. GitHub: CI is connected and green; no deployment credential is configured or needed.

Never paste credentials into chat. Put them in a local `.env` or deployment secret/environment settings.
