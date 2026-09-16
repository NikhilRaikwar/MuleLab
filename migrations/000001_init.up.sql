CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE agent_status AS ENUM ('PROPOSED','POLICY_CHECKED','AWAITING_APPROVAL','APPROVED','RUNNING','EVALUATING','WINNING','LOSING','PAUSED','RETIRED','FAILED');
CREATE TYPE agent_type AS ENUM ('SUPERVISOR','STORE','NOTIFY','GIVE');

CREATE TABLE businesses (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slug text NOT NULL UNIQUE,
  name text NOT NULL,
  description text NOT NULL,
  scenario_seed bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE business_snapshots (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id uuid NOT NULL REFERENCES businesses(id),
  day integer NOT NULL CHECK (day >= 0),
  metrics jsonb NOT NULL CHECK (jsonb_typeof(metrics) = 'object'),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (business_id, day)
);

CREATE TABLE objectives (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id uuid NOT NULL REFERENCES businesses(id),
  kind text NOT NULL,
  description text NOT NULL,
  target_value numeric NOT NULL,
  budget_usd numeric(10,2) NOT NULL CHECK (budget_usd >= 0),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE agent_runs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  trace_id uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  business_id uuid NOT NULL REFERENCES businesses(id),
  objective_id uuid NOT NULL REFERENCES objectives(id),
  status agent_status NOT NULL DEFAULT 'PROPOSED',
  scenario_seed bigint NOT NULL,
  day integer NOT NULL DEFAULT 0 CHECK (day >= 0),
  model_call_budget integer NOT NULL CHECK (model_call_budget BETWEEN 0 AND 8),
  tool_call_budget integer NOT NULL CHECK (tool_call_budget BETWEEN 0 AND 12),
  deadline_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE agent_instances (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id uuid NOT NULL REFERENCES agent_runs(id) ON DELETE CASCADE,
  type agent_type NOT NULL,
  status agent_status NOT NULL DEFAULT 'PROPOSED',
  current_task text NOT NULL DEFAULT '',
  policy_violation_count integer NOT NULL DEFAULT 0 CHECK (policy_violation_count >= 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (run_id, type)
);

CREATE TABLE agent_proposals (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id uuid NOT NULL REFERENCES agent_runs(id) ON DELETE CASCADE,
  agent_instance_id uuid NOT NULL REFERENCES agent_instances(id),
  version integer NOT NULL CHECK (version > 0),
  payload jsonb NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
  schema_valid boolean NOT NULL,
  policy_allowed boolean NOT NULL,
  requires_approval boolean NOT NULL,
  immutable_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (agent_instance_id, version)
);

CREATE TABLE experiments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id uuid NOT NULL REFERENCES agent_runs(id) ON DELETE CASCADE,
  proposal_id uuid NOT NULL UNIQUE REFERENCES agent_proposals(id),
  status agent_status NOT NULL,
  budget_usd numeric(10,2) NOT NULL CHECK (budget_usd >= 0),
  primary_metric text NOT NULL,
  success_threshold jsonb NOT NULL,
  stop_condition jsonb NOT NULL,
  approved_at timestamptz,
  started_at timestamptz,
  completed_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE experiment_results (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  experiment_id uuid NOT NULL UNIQUE REFERENCES experiments(id) ON DELETE CASCADE,
  sample_size integer NOT NULL CHECK (sample_size >= 0),
  baseline numeric NOT NULL,
  observed numeric NOT NULL,
  spend_usd numeric(10,2) NOT NULL CHECK (spend_usd >= 0),
  successful boolean NOT NULL,
  stop_condition_met boolean NOT NULL,
  calculated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE human_approvals (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  proposal_id uuid NOT NULL REFERENCES agent_proposals(id),
  proposal_version integer NOT NULL,
  decision text NOT NULL CHECK (decision IN ('APPROVED','REJECTED')),
  exact_arguments_hash text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (proposal_id, proposal_version)
);

CREATE TABLE tool_calls (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id uuid NOT NULL REFERENCES agent_runs(id) ON DELETE CASCADE,
  experiment_id uuid REFERENCES experiments(id),
  tool_name text NOT NULL,
  idempotency_key text NOT NULL UNIQUE,
  input jsonb NOT NULL,
  output jsonb,
  status text NOT NULL CHECK (status IN ('STARTED','SUCCEEDED','FAILED','TIMED_OUT')),
  error_code text,
  latency_ms integer CHECK (latency_ms >= 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz
);

CREATE TABLE model_profiles (
  id text PRIMARY KEY,
  supports_json_schema boolean NOT NULL,
  supports_tools boolean NOT NULL,
  is_free_allowed boolean NOT NULL,
  purpose text NOT NULL,
  refreshed_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE model_calls (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id uuid REFERENCES agent_runs(id) ON DELETE CASCADE,
  requested_model text NOT NULL,
  resolved_model text NOT NULL,
  prompt_version text NOT NULL,
  structured_output_valid boolean NOT NULL,
  prompt_tokens integer CHECK (prompt_tokens >= 0),
  completion_tokens integer CHECK (completion_tokens >= 0),
  cost_usd numeric(14,8) CHECK (cost_usd >= 0),
  latency_ms integer NOT NULL CHECK (latency_ms >= 0),
  retry_count integer NOT NULL DEFAULT 0 CHECK (retry_count BETWEEN 0 AND 1),
  failure_code text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE trace_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  trace_id uuid NOT NULL,
  run_id uuid REFERENCES agent_runs(id) ON DELETE CASCADE,
  sequence integer NOT NULL CHECK (sequence > 0),
  category text NOT NULL,
  title text NOT NULL,
  summary text NOT NULL,
  evidence jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (trace_id, sequence)
);

CREATE TABLE policies (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  version integer NOT NULL CHECK (version > 0),
  rules jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (name, version)
);

CREATE TABLE eval_cases (
  id text PRIMARY KEY,
  name text NOT NULL,
  category text NOT NULL,
  fixture jsonb NOT NULL,
  expected jsonb NOT NULL,
  release_blocking boolean NOT NULL DEFAULT false
);

CREATE TABLE eval_runs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  model_id text NOT NULL,
  source text NOT NULL CHECK (source IN ('FIXTURE','DETERMINISTIC_FALLBACK','LIVE')),
  passed boolean NOT NULL,
  release_blocked boolean NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE eval_results (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  eval_run_id uuid NOT NULL REFERENCES eval_runs(id) ON DELETE CASCADE,
  eval_case_id text NOT NULL REFERENCES eval_cases(id),
  passed boolean NOT NULL,
  deterministic boolean NOT NULL,
  score numeric NOT NULL CHECK (score BETWEEN 0 AND 1),
  details text NOT NULL,
  UNIQUE (eval_run_id, eval_case_id)
);

CREATE INDEX agent_runs_business_created_idx ON agent_runs (business_id, created_at DESC);
CREATE INDEX agent_instances_run_status_idx ON agent_instances (run_id, status);
CREATE INDEX experiments_run_status_idx ON experiments (run_id, status);
CREATE INDEX tool_calls_run_created_idx ON tool_calls (run_id, created_at);
CREATE INDEX model_calls_run_created_idx ON model_calls (run_id, created_at);
CREATE INDEX trace_events_run_sequence_idx ON trace_events (run_id, sequence);
CREATE INDEX eval_results_run_idx ON eval_results (eval_run_id);

