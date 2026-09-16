# Environment contract

Copy `.env.example` to `.env`. `.env` is Git-ignored and is the only local file that may contain secrets.

| Variable | Required | Secret | Scope | Notes |
| --- | --- | --- | --- | --- |
| `APP_ENV`, `APP_URL`, `API_PUBLIC_URL`, `PORT`, `CORS_ALLOWED_ORIGINS`, `LOG_LEVEL` | No | No | Local/production | Service address and log behavior. |
| `DATABASE_URL` | Yes for API | Yes for Neon; local Docker only is not sensitive | Local/production | Pooled Neon URL in hosted runtime; `postgres` Docker URL locally. |
| `DATABASE_URL_DIRECT` | Migration only | Yes | Production | Direct, non-pooler Neon URL. Migrations must not use the pooled URL. |
| `OPENROUTER_API_KEY` | No | Yes | Local/production | Server-side only. Empty activates the clearly labelled deterministic planner path. |
| `OPENROUTER_BASE_URL`, `OPENROUTER_MODEL`, `OPENROUTER_FALLBACK_MODEL`, `OPENROUTER_HTTP_REFERER`, `OPENROUTER_APP_NAME` | No | No | Local/production | Defaults to the free OpenRouter route. Any paid model is rejected when zero-cost mode is enabled. |
| `MODEL_REQUEST_TIMEOUT_SECONDS`, `MODEL_MAX_RETRIES`, `MODEL_MAX_OUTPUT_TOKENS`, `MAX_MODEL_CALLS_PER_RUN` | No | No | Local/production | Bounded model-call safety controls. |
| `ZERO_COST_MODE`, `ALLOW_PAID_MODELS`, `DETERMINISTIC_MODEL_FALLBACK` | No | No | Local/production | Keep `ZERO_COST_MODE=true` and `ALLOW_PAID_MODELS=false` for P0. |
| `MAX_AGENT_ITERATIONS`, `MAX_TOOL_CALLS_PER_RUN`, `MAX_RUN_DURATION_SECONDS`, `MAX_DAILY_PUBLIC_RUNS` | No | No | Local/production | Agent/run hard limits. |
| `MAX_EXPERIMENT_BUDGET_USD`, approval and messaging/discount variables | No | No | Local/production | Deterministic authority limits. |
| `OTEL_*`, `TRACE_REDACTION_ENABLED` | No | OTLP headers may be secret | Optional observability | Product traces persist regardless of OTLP export. |
| `NEXT_PUBLIC_*` | No | No | Browser-visible | Never add an OpenRouter or database credential to a `NEXT_PUBLIC_` value. |
| `GCP_*`, `MAX_CLOUD_RUN_INSTANCES` | Deployment only | No | Production | Project configuration only; credentials remain outside the repository. |

## Local Docker values

After Docker is available, put this only in your local `.env`:

```dotenv
DATABASE_URL=postgresql://mulelab:mulelab@localhost:5432/mulelab?sslmode=disable
DATABASE_URL_DIRECT=postgresql://mulelab:mulelab@localhost:5432/mulelab?sslmode=disable
```

## Neon values

Use `DATABASE_URL` for pooled runtime traffic and `DATABASE_URL_DIRECT` for migrations. Neon pooled hosts contain `-pooler`; the direct migration host does not. Never place either real URL in committed files.
