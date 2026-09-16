package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/raikwar/mulelab/apps/api/graph/model"
	"github.com/raikwar/mulelab/internal/domain"
	"github.com/raikwar/mulelab/internal/runtime"
)

func (r *mutationResolver) CreateRun(ctx context.Context, input model.CreateRunInput) (*model.AgentRun, error) {
	target, budget := 50.0, 100.0
	if input.TargetValue != nil {
		target = *input.TargetValue
	}
	if input.BudgetUsd != nil {
		budget = *input.BudgetUsd
	}
	seed := 0
	if input.Seed != nil {
		seed = *input.Seed
	}
	live := false
	if input.LiveModel != nil {
		live = *input.LiveModel
	}
	run, err := r.Runtime.CreateRun(ctx, input.BusinessID, input.ObjectiveKind, target, budget, seed, live)
	if err != nil {
		return nil, gqlError("CREATE_RUN_FAILED", err)
	}
	return mapRun(run), nil
}
func (r *mutationResolver) GenerateProposals(ctx context.Context, runID string) ([]*model.ExperimentProposal, error) {
	items, err := r.Runtime.GenerateProposals(ctx, runID)
	if err != nil {
		return nil, gqlError("PROPOSAL_FAILED", err)
	}
	out := make([]*model.ExperimentProposal, len(items))
	for i, p := range items {
		out[i] = mapProposal(p)
	}
	return out, nil
}
func (r *mutationResolver) ApproveExperiment(ctx context.Context, id string) (*model.Experiment, error) {
	exp, err := r.Runtime.Approve(ctx, id)
	if err != nil {
		return nil, gqlError("APPROVAL_FAILED", err)
	}
	run, err := r.Runtime.RunForProposal(id)
	if err != nil {
		return nil, gqlError("RUN_NOT_FOUND", err)
	}
	return mapExperiment(run, exp), nil
}
func (r *mutationResolver) RejectExperiment(ctx context.Context, id string) (*model.Experiment, error) {
	exp, err := r.Runtime.Reject(ctx, id)
	if err != nil {
		return nil, gqlError("REJECTION_FAILED", err)
	}
	run, err := r.Runtime.RunForProposal(id)
	if err != nil {
		return nil, gqlError("RUN_NOT_FOUND", err)
	}
	return mapExperiment(run, exp), nil
}
func (r *mutationResolver) AdvanceSimulation(ctx context.Context, runID string, days int) (*model.AgentRun, error) {
	run, err := r.Runtime.Advance(ctx, runID, days)
	if err != nil {
		return nil, gqlError("SIMULATION_FAILED", err)
	}
	return mapRun(run), nil
}
func (r *mutationResolver) EvaluateRun(ctx context.Context, runID string) (*model.AgentRun, error) {
	run, err := r.Runtime.GetRun(ctx, runID)
	if err != nil {
		return nil, gqlError("RUN_NOT_FOUND", err)
	}
	return mapRun(run), nil
}
func (r *mutationResolver) RunEvalSuite(ctx context.Context) (*model.EvalReport, error) {
	report, err := r.Runtime.RunEvalSuite(ctx)
	if err != nil {
		return nil, gqlError("EVAL_RUN_FAILED", err)
	}
	return mapEvalReport(report), nil
}
func (r *mutationResolver) ReplayScenario(ctx context.Context, input model.ReplayInput) (*model.ReplayResult, error) {
	run, err := r.Runtime.GetRun(ctx, input.RunID)
	if err != nil {
		return nil, gqlError("RUN_NOT_FOUND", err)
	}
	items := []*model.ModelComparison{}
	for _, id := range input.ModelIds {
		if id == "deterministic/demo-v1" {
			score := 1.0
			zero := 0
			latency := 0
			tokens := 0
			cost := 0.0
			items = append(items, &model.ModelComparison{ModelID: id, Executed: true, Source: "FIXTURE", Score: &score, PolicyViolations: &zero, ToolCallSuccessRate: &score, StructuredOutputFailures: &zero, LatencyMs: &latency, TotalTokens: &tokens, EstimatedCostUsd: &cost, Retries: &zero, Note: "Fixture-backed deterministic baseline."})
		} else {
			items = append(items, &model.ModelComparison{ModelID: id, Executed: false, Source: "NOT_RUN", Note: "Only configured free models may be run; no comparison was fabricated."})
		}
	}
	return &model.ReplayResult{ScenarioSeed: int(run.Seed), Comparisons: items}, nil
}
func (r *queryResolver) Health(ctx context.Context) (string, error) { return "ok", nil }
func (r *queryResolver) DemoBusinesses(ctx context.Context) ([]*model.Business, error) {
	items, err := r.Runtime.Businesses(ctx)
	if err != nil {
		return nil, gqlError("DATABASE_UNAVAILABLE", err)
	}
	out := make([]*model.Business, len(items))
	for i, b := range items {
		out[i] = mapBusiness(b)
	}
	return out, nil
}
func (r *queryResolver) Business(ctx context.Context, id string) (*model.Business, error) {
	b, err := r.Runtime.Business(ctx, id)
	if err != nil {
		return nil, nil
	}
	return mapBusiness(b), nil
}
func (r *queryResolver) Run(ctx context.Context, id string) (*model.AgentRun, error) {
	run, err := r.Runtime.GetRun(ctx, id)
	if err != nil {
		return nil, nil
	}
	return mapRun(run), nil
}
func (r *queryResolver) Traces(ctx context.Context, runID string) ([]*model.TraceEvent, error) {
	items, err := r.Runtime.Traces(ctx, runID)
	if err != nil {
		return nil, gqlError("RUN_NOT_FOUND", err)
	}
	out := make([]*model.TraceEvent, len(items))
	for i, t := range items {
		out[i] = mapTrace(t)
	}
	return out, nil
}
func (r *queryResolver) LatestEvalReport(ctx context.Context) (*model.EvalReport, error) {
	report, err := r.Runtime.LatestEvalReport(ctx)
	if err != nil {
		return nil, gqlError("EVAL_REPORT_NOT_FOUND", err)
	}
	return mapEvalReport(report), nil
}
func (r *queryResolver) ModelProfiles(ctx context.Context) ([]*model.ModelProfile, error) {
	return []*model.ModelProfile{{ID: "openrouter/free", SupportsJSONSchema: true, SupportsTools: true, IsFreeAllowed: true, Purpose: "Dynamic free route; resolved model is persisted when used."}, {ID: "deterministic/demo-v1", SupportsJSONSchema: true, SupportsTools: true, IsFreeAllowed: true, Purpose: "Labeled safe demo fallback and CI fixture."}}, nil
}
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func mapBusiness(b domain.Business) *model.Business {
	return &model.Business{ID: b.ID, Name: b.Name, Description: b.Description, Seed: int(b.Seed), Metrics: mapMetrics(b.Metrics)}
}
func mapMetrics(m domain.Metrics) *model.BusinessMetrics {
	return &model.BusinessMetrics{Orders: m.Orders, ConversionRate: m.ConversionRate, RepeatPurchaseRate: m.RepeatPurchaseRate, Subscribers: m.Subscribers, QualifiedSubscribers: m.QualifiedSubscribers, ExperimentSpendUsd: m.ExperimentSpendUSD}
}
func mapRun(run *runtime.Run) *model.AgentRun {
	agents := make([]*model.AgentInstance, len(run.Agents))
	for i, a := range run.Agents {
		agents[i] = &model.AgentInstance{ID: a.ID, Type: model.AgentType(a.Type), Name: a.Name, Status: model.AgentStatus(a.Status), CurrentTask: a.CurrentTask, CumulativeImpact: a.CumulativeImpact, PolicyViolations: a.PolicyViolationCount}
	}
	props := make([]*model.ExperimentProposal, len(run.Proposals))
	for i, p := range run.Proposals {
		props[i] = mapProposal(p)
	}
	exps := make([]*model.Experiment, len(run.Experiments))
	for i, e := range run.Experiments {
		exps[i] = mapExperiment(run, e)
	}
	stats := make([]*model.ModelRunStats, len(run.ModelStats))
	for i, s := range run.ModelStats {
		stats[i] = &model.ModelRunStats{RequestedModel: s.RequestedModel, ResolvedModel: s.ResolvedModel, Source: s.Source, LatencyMs: int(s.Latency.Milliseconds()), PromptTokens: s.PromptTokens, CompletionTokens: s.CompletionTokens, EstimatedCostUsd: s.CostUSD, Retries: s.Retries, StructuredOutputValid: true}
	}
	return &model.AgentRun{ID: run.ID, TraceID: run.TraceID, BusinessID: run.Business.ID, Objective: &model.Objective{ID: run.Objective.ID, Kind: run.Objective.Kind, Description: run.Objective.Description, TargetValue: run.Objective.TargetValue, BudgetUsd: run.Objective.BudgetUSD}, Status: model.AgentStatus(run.Status), Day: run.Day, Seed: int(run.Seed), Agents: agents, Proposals: props, Experiments: exps, ModelStats: stats, Metrics: mapMetrics(run.Business.Metrics), CreatedAt: run.CreatedAt, UpdatedAt: run.UpdatedAt}
}
func mapProposal(p domain.Proposal) *model.ExperimentProposal {
	raw, _ := json.Marshal(p.ToolArguments)
	return &model.ExperimentProposal{ID: p.ID, Version: p.Version, AgentType: model.AgentType(p.AgentType), Hypothesis: p.Hypothesis, ToolName: p.ToolName, ToolArgumentsJSON: string(raw), PrimaryMetric: p.PrimaryMetric, ExpectedOutcome: p.ExpectedOutcome, BudgetUsd: p.BudgetUSD, SuccessThreshold: p.SuccessThreshold, StopCondition: p.StopCondition, EvidenceIds: p.EvidenceIDs, Confidence: p.Confidence, Status: model.AgentStatus(p.Status), RequiresApproval: p.RequiresApproval, PolicyReasons: p.PolicyReasons}
}
func mapExperiment(run *runtime.Run, e domain.Experiment) *model.Experiment {
	var p domain.Proposal
	for _, candidate := range run.Proposals {
		if candidate.ID == e.ProposalID {
			p = candidate
			break
		}
	}
	out := &model.Experiment{ID: e.ID, Proposal: mapProposal(p), Status: model.AgentStatus(e.Status), ApprovedAt: e.ApprovedAt}
	if e.Result != nil {
		out.Result = &model.ExperimentResult{PrimaryMetric: e.Result.PrimaryMetric, Baseline: e.Result.Baseline, Observed: e.Result.Observed, SampleSize: e.Result.SampleSize, SpendUsd: e.Result.SpendUSD, Successful: e.Result.Successful, StopConditionMet: e.Result.StopConditionMet, Summary: e.Result.Summary}
	}
	return out
}
func mapTrace(t domain.TraceEvent) *model.TraceEvent {
	raw, _ := json.Marshal(t.Evidence)
	return &model.TraceEvent{ID: t.ID, TraceID: t.TraceID, Sequence: t.Sequence, Category: t.Category, Title: t.Title, Summary: t.Summary, EvidenceJSON: string(raw), CreatedAt: t.CreatedAt}
}
func mapEvalReport(report domain.EvalReport) *model.EvalReport {
	cases := make([]*model.EvalCaseResult, len(report.Cases))
	passed := 0
	for i, c := range report.Cases {
		if c.Passed {
			passed++
		}
		cases[i] = &model.EvalCaseResult{ID: c.ID, Name: c.Name, Category: c.Category, Passed: c.Passed, Deterministic: c.Deterministic, Score: c.Score, Details: c.Details}
	}
	return &model.EvalReport{ID: report.ID, Passed: report.Passed, ReleaseBlocked: report.ReleaseBlocked, Total: len(cases), PassedCount: passed, Cases: cases, CreatedAt: report.CreatedAt}
}
func gqlError(code string, err error) error { return fmt.Errorf("%s: %w", code, err) }
