# Local setup

## Prerequisites

- Go 1.27.x
- Node.js 20.9+ (Node 22 LTS is suitable)
- npm
- PostgreSQL 17+ or Docker Desktop / Docker Compose

This project has been verified with Go 1.27, Node 22, Docker Desktop, Docker Compose, and PostgreSQL 17. The local Compose path uses PostgreSQL 17; hosted development has also been verified against Neon PostgreSQL 17. Use `D:\DevTemp` for Windows frontend worker temporary files when the system temp directory is not writable.

## Configuration

Copy `.env.example` to `.env` and `apps/web/.env.local.example` to `apps/web/.env.local`. Keep both local files uncommitted.

The deterministic demo does not require an OpenRouter key. Live AI does.

## Secrets Nikhil supplies

### `OPENROUTER_API_KEY`

- Required: live AI only
- Secret: yes
- Cost: free route by default; paid models are blocked in zero-cost mode
- Obtain: OpenRouter dashboard → Keys
- Format: `sk-or-v1-...`
- Store: local `.env` / Cloud Run environment or secret; never commit

### `DATABASE_URL` and `DATABASE_URL_DIRECT`

- Required: database-backed local/native or hosted persistence
- Secret: yes
- Cost: use Neon Free
- Obtain: Neon project → Connect → pooled connection string (`DATABASE_URL`) and direct migration URL (`DATABASE_URL_DIRECT`)
- Format: `postgresql://user:password@ep-...-pooler.../neondb?sslmode=require`
- Store: local `.env` / future Cloud Run configuration; never commit

### `GCP_PROJECT_ID` and gcloud authentication

- Required: deployment only
- Secret: project ID no; credentials yes
- Cost: GCP resources can incur charges
- Obtain: Google Cloud project selector
- Store: project ID may be an environment value; authenticate with Google Cloud CLI, never commit service-account keys
