# Deploy to Google Cloud Run

No deployment command should be executed until the explicit cost check is reviewed and approved.

## DEPLOYMENT COST CHECK

| Resource | Expected free-tier behavior | Possible charge trigger | Protection |
| --- | --- | --- | --- |
| Cloud Run API | Scale-to-zero and small free allowance | Traffic, CPU, memory, egress beyond allowance | min 0, max 1, conservative limits, public run cap |
| Cloud Run web | Same | Same | min 0, max 1, static-heavy UI |
| Artifact Registry | Small storage allowance may apply | Stored image volume and egress | Small multi-stage images; limited retained tags |
| Cloud Build | Build-minute allowance may apply | Repeated/long builds | Build locally or intentionally; no automatic deploy yet |
| Neon | Free project | Provider plan limits/upgrade | Free plan; small seed/trace data |
| OpenRouter | Free models | Selecting paid model or exceeding policy | `ZERO_COST_MODE=true`, `ALLOW_PAID_MODELS=false` |
| Cloud Logging | Included ingestion allowance | Excessive logs/retention | structured concise logs; no payload/secret logging |

## Required authorization

Authenticate with Google Cloud CLI, select the intended billing-enabled project, and explicitly approve resource creation. Never commit service-account JSON.

## Planned commands

```sh
gcloud auth login
gcloud config set project "$GCP_PROJECT_ID"
gcloud services enable run.googleapis.com artifactregistry.googleapis.com cloudbuild.googleapis.com
gcloud artifacts repositories create "$GCP_ARTIFACT_REPOSITORY" --repository-format=docker --location="$GCP_REGION"
gcloud auth configure-docker "$GCP_REGION-docker.pkg.dev"

docker build -f apps/api/Dockerfile -t "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_ARTIFACT_REPOSITORY/mulelab-api:latest" .
docker push "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_ARTIFACT_REPOSITORY/mulelab-api:latest"
gcloud run deploy "$GCP_API_SERVICE_NAME" --image "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_ARTIFACT_REPOSITORY/mulelab-api:latest" --region "$GCP_REGION" --allow-unauthenticated --min 0 --max 1 --cpu 1 --memory 512Mi --set-env-vars "APP_ENV=production,ZERO_COST_MODE=true,MAX_DAILY_PUBLIC_RUNS=50"

docker build -f apps/web/Dockerfile -t "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_ARTIFACT_REPOSITORY/mulelab-web:latest" .
docker push "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_ARTIFACT_REPOSITORY/mulelab-web:latest"
gcloud run deploy "$GCP_WEB_SERVICE_NAME" --image "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_ARTIFACT_REPOSITORY/mulelab-web:latest" --region "$GCP_REGION" --allow-unauthenticated --min 0 --max 1 --cpu 1 --memory 512Mi
```

Secrets and the final API URL must be configured through Cloud Run environment/secret settings after the services exist. Exact commands will be finalized and tested before deployment.
