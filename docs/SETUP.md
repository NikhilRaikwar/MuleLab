# Local setup

## Prerequisites

- Go 1.27.x
- Node.js 20.9+ (Node 22 LTS is suitable)
- npm
- PostgreSQL 16+ or Docker

This host currently has Node 22.17.0, npm 10.9.2, Git 2.50.1, and Go 1.27.0. Backend tests and frontend lint/typecheck/production build have passed. Docker Desktop is not yet available to this shell, so container verification requires either restarting the shell after Docker installation or correcting Docker's PATH/service state.

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

### `DATABASE_URL`

- Required: deployed persistence
- Secret: yes
- Cost: use Neon Free
- Obtain: Neon project → Connect → pooled connection string
- Format: `postgresql://user:password@ep-...-pooler.../neondb?sslmode=require`
- Store: local `.env` / Cloud Run configuration; never commit

### `GCP_PROJECT_ID` and gcloud authentication

- Required: deployment only
- Secret: project ID no; credentials yes
- Cost: GCP resources can incur charges
- Obtain: Google Cloud project selector
- Store: project ID may be an environment value; authenticate with Google Cloud CLI, never commit service-account keys
