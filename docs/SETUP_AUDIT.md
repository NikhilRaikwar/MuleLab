# Setup audit

Audited 2026-09-16. The repository initially contained only `MuleLab_PRD.md`; there was no source, Git metadata, build configuration, CI, database schema, deployment configuration, or documentation.

| Service | Why needed | Required for MVP? | Connection/plugin needed? | Secret/env required? | Free option? | Status |
| --- | --- | --- | --- | --- | --- | --- |
| GitHub | Source hosting and CI | Required for hosted CI, not local build | Git remote and authenticated `gh`/Git | No app secret | Public repo and Actions allowance | Not connected |
| OpenRouter | Live structured agent proposals and arena | Optional for deterministic public demo; required for live AI | OpenRouter account/API key; no Codex plugin | `OPENROUTER_API_KEY` | `openrouter/free`, rate-limited | Key configured locally; metadata verified without inference |
| GCP | Public Cloud Run deployment | Required for target deployment | `gcloud` login and authorized project | Local auth; project/region env | Cloud Run has free allowance, never guaranteed $0 | Not connected |
| Neon Postgres | Deployed durable state | Required for deployed persistence | Neon account/project and Neon CLI | `DATABASE_URL`, `DATABASE_URL_DIRECT` | Neon Free | Blocked: authenticated organization is Vercel-managed and rejects project creation via Neon CLI |
| Local PostgreSQL | Local integration tests | Recommended | Docker or local Postgres | Local-only URL | Yes | Docker unavailable on this host |
| Artifact Registry | Store deployment images | Required for image-based Cloud Run deploy | GCP authentication | No repository credential committed | Small free allowance may apply | Not connected |
| OpenTelemetry collector/vendor | External trace export | No | Optional OTLP endpoint | Endpoint/headers | Local/stdout and several free options | Disabled |
| GitHub Actions | Deterministic CI | Required once hosted | GitHub repository | Optional deployment secrets | Included allowance varies | Not configured yet |

## Local tooling observation

Go 1.27.0 is installed. gqlgen generation succeeds using an isolated verified module cache. Focused backend tests, frontend lint/typecheck, and the production frontend build pass. Docker Desktop is still not accessible from this shell (`docker` is not on `PATH`), so container and local-Postgres verification remain pending.

Neon CLI authentication was verified. The selected organization had no existing projects, but Neon rejected creating the explicitly approved `mulelab` project because the organization is managed by Vercel. No project, branch, paid resource, or credential was created. A standard Neon organization (or a Vercel-managed project created through Vercel itself) is required before hosted-Neon verification can proceed.

## Connection actions

No Codex plugin or MCP server is required. Development proceeds with deterministic adapters and fixtures.

1. OpenRouter: the local key enables a bounded optional live summary call. The public demo remains usable without it via a labeled deterministic fallback.
2. Neon: use or create a non-Vercel-managed Neon organization, then create a separate MuleLab Free project and copy its pooled and direct connection strings. Local work can proceed without it.
3. GCP: authenticate `gcloud` and select a billing-enabled project only after the deployment cost review. No cloud resources will be created without explicit approval.
4. GitHub: create/connect a repository when CI and public source hosting are desired. Local implementation does not depend on it.

Never paste credentials into chat. Put them in a local `.env` or deployment secret/environment settings.
