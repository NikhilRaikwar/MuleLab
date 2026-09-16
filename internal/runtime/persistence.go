package runtime

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/raikwar/mulelab/internal/domain"
	"github.com/raikwar/mulelab/internal/modelgateway"
	"github.com/raikwar/mulelab/internal/tools"
)

func (r *Runtime) persistRun(ctx context.Context, run *Run) error {
	var businessUUID string
	if err := r.db.Pool.QueryRow(ctx, `SELECT id FROM businesses WHERE slug=$1`, run.Business.ID).Scan(&businessUUID); err != nil {
		return err
	}
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO objectives(id,business_id,kind,description,target_value,budget_usd) VALUES($1,$2,$3,$4,$5,$6)`, run.Objective.ID, businessUUID, run.Objective.Kind, run.Objective.Description, run.Objective.TargetValue, run.Objective.BudgetUSD)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO agent_runs(id,trace_id,business_id,objective_id,status,scenario_seed,model_call_budget,tool_call_budget,deadline_at) VALUES($1,$2,$3,$4,$5::agent_status,$6,8,12,$7)`, run.ID, run.TraceID, businessUUID, run.Objective.ID, string(run.Status), run.Seed, time.Now().Add(90*time.Second))
	if err != nil {
		return err
	}
	for _, a := range run.Agents {
		if _, err = tx.Exec(ctx, `INSERT INTO agent_instances(id,run_id,type,status,current_task) VALUES($1,$2,$3::agent_type,$4::agent_status,$5)`, a.ID, run.ID, string(a.Type), string(a.Status), a.CurrentTask); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
func (r *Runtime) updateRun(ctx context.Context, run *Run) {
	_, _ = r.db.Pool.Exec(ctx, `UPDATE agent_runs SET status=$2::agent_status,day=$3,updated_at=now() WHERE id=$1`, run.ID, string(run.Status), run.Day)
}
func (r *Runtime) persistProposal(ctx context.Context, run *Run, p domain.Proposal) {
	payload, _ := json.Marshal(p)
	var agentID string
	_ = r.db.Pool.QueryRow(ctx, `SELECT id FROM agent_instances WHERE run_id=$1 AND type=$2::agent_type`, run.ID, string(p.AgentType)).Scan(&agentID)
	_, _ = r.db.Pool.Exec(ctx, `INSERT INTO agent_proposals(id,run_id,agent_instance_id,version,payload,schema_valid,policy_allowed,requires_approval) VALUES($1,$2,$3,$4,$5::jsonb,true,true,$6)`, p.ID, run.ID, agentID, p.Version, payload, p.RequiresApproval)
}
func thresholdJSON(value string) string {
	b, _ := json.Marshal(map[string]string{"expression": value})
	return string(b)
}
func (r *Runtime) persistApproval(ctx context.Context, run *Run, p domain.Proposal, exp domain.Experiment) {
	args, _ := json.Marshal(p.ToolArguments)
	sum := sha256.Sum256(args)
	_, _ = r.db.Pool.Exec(ctx, `INSERT INTO human_approvals(proposal_id,proposal_version,decision,exact_arguments_hash) VALUES($1,$2,'APPROVED',$3) ON CONFLICT(proposal_id,proposal_version) DO NOTHING`, p.ID, p.Version, hex.EncodeToString(sum[:]))
	_, _ = r.db.Pool.Exec(ctx, `UPDATE agent_proposals SET immutable_at=now() WHERE id=$1`, p.ID)
	_, _ = r.db.Pool.Exec(ctx, `INSERT INTO experiments(id,run_id,proposal_id,status,budget_usd,primary_metric,success_threshold,stop_condition,approved_at) VALUES($1,$2,$3,'APPROVED'::agent_status,$4,$5,$6::jsonb,$7::jsonb,now()) ON CONFLICT (proposal_id) DO UPDATE SET status='APPROVED'::agent_status,approved_at=now()`, exp.ID, run.ID, p.ID, p.BudgetUSD, p.PrimaryMetric, thresholdJSON(p.SuccessThreshold), thresholdJSON(p.StopCondition))
}
func (r *Runtime) persistRejection(ctx context.Context, run *Run, p domain.Proposal, exp domain.Experiment) error {
	args, err := json.Marshal(p.ToolArguments)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(args)
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `INSERT INTO human_approvals(proposal_id,proposal_version,decision,exact_arguments_hash) VALUES($1,$2,'REJECTED',$3) ON CONFLICT(proposal_id,proposal_version) DO NOTHING`, p.ID, p.Version, hex.EncodeToString(sum[:])); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO experiments(id,run_id,proposal_id,status,budget_usd,primary_metric,success_threshold,stop_condition) VALUES($1,$2,$3,'PAUSED'::agent_status,$4,$5,$6::jsonb,$7::jsonb) ON CONFLICT (proposal_id) DO UPDATE SET status='PAUSED'::agent_status`, exp.ID, run.ID, p.ID, p.BudgetUSD, p.PrimaryMetric, thresholdJSON(p.SuccessThreshold), thresholdJSON(p.StopCondition))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Runtime) persistResult(ctx context.Context, run *Run, exp domain.Experiment, p domain.Proposal) {
	if exp.Result == nil {
		return
	}
	result := exp.Result
	_, _ = r.db.Pool.Exec(ctx, `INSERT INTO experiment_results(experiment_id,sample_size,baseline,observed,spend_usd,successful,stop_condition_met) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(experiment_id) DO UPDATE SET sample_size=EXCLUDED.sample_size,baseline=EXCLUDED.baseline,observed=EXCLUDED.observed,spend_usd=EXCLUDED.spend_usd,successful=EXCLUDED.successful,stop_condition_met=EXCLUDED.stop_condition_met`, exp.ID, result.SampleSize, result.Baseline, result.Observed, result.SpendUSD, result.Successful, result.StopConditionMet)
	_, _ = r.db.Pool.Exec(ctx, `UPDATE experiments SET status=$2::agent_status,completed_at=now() WHERE id=$1`, exp.ID, string(exp.Status))
	var agentID string
	_ = r.db.Pool.QueryRow(ctx, `SELECT id FROM agent_instances WHERE run_id=$1 AND type=$2::agent_type`, run.ID, string(p.AgentType)).Scan(&agentID)
	_, _ = r.db.Pool.Exec(ctx, `UPDATE agent_instances SET status=$2::agent_status,updated_at=now() WHERE id=$1`, agentID, string(exp.Status))
}
func (r *Runtime) persistModelCall(ctx context.Context, run *Run, result modelgateway.Result) {
	_, _ = r.db.Pool.Exec(ctx, `INSERT INTO model_calls(run_id,requested_model,resolved_model,prompt_version,structured_output_valid,prompt_tokens,completion_tokens,cost_usd,latency_ms,retry_count) VALUES($1,$2,$3,'supervisor-summary-v1',true,$4,$5,$6,$7,$8)`, run.ID, result.RequestedModel, result.ResolvedModel, result.PromptTokens, result.CompletionTokens, result.CostUSD, result.Latency.Milliseconds(), result.Retries)
}
func (r *Runtime) persistToolCall(ctx context.Context, run *Run, exp domain.Experiment, call tools.Call) error {
	input, err := json.Marshal(call.Proposal.ToolArguments)
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, `INSERT INTO tool_calls(run_id,experiment_id,tool_name,idempotency_key,input,output,status,latency_ms,completed_at) VALUES($1,$2,$3,$4,$5::jsonb,'{}'::jsonb,'SUCCEEDED',0,now()) ON CONFLICT (idempotency_key) DO NOTHING`, run.ID, exp.ID, call.Name, call.IdempotencyKey, input)
	return err
}
func (r *Runtime) trace(ctx context.Context, run *Run, category, title, summary string, evidence map[string]any) {
	e := domain.TraceEvent{ID: uuid.NewString(), TraceID: run.TraceID, Sequence: len(run.Traces) + 1, Category: category, Title: title, Summary: summary, Evidence: evidence, CreatedAt: time.Now().UTC()}
	run.Traces = append(run.Traces, e)
	raw, _ := json.Marshal(evidence)
	_, _ = r.db.Pool.Exec(ctx, `INSERT INTO trace_events(trace_id,run_id,sequence,category,title,summary,evidence) VALUES($1,$2,$3,$4,$5,$6,$7::jsonb)`, e.TraceID, run.ID, e.Sequence, e.Category, e.Title, e.Summary, raw)
}
