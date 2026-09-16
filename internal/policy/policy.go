package policy

import (
	"fmt"
	"strings"

	"github.com/raikwar/mulelab/internal/domain"
)

type Limits struct {
	RunBudgetUSD               float64
	ApprovalSpendUSD           float64
	ApprovalAudience           int
	ApprovalDiscountPct        float64
	ApprovalConfidence         float64
	MaxDiscountPct             float64
	MaxMessagesPerContactPer7d int
}

type Decision struct {
	Allowed          bool
	RequiresApproval bool
	Reasons          []string
}

var permissions = map[domain.AgentType]map[string]bool{
	domain.AgentStore:  {"store.get_summary": true, "store.list_products": true, "store.get_product_metrics": true, "store.create_product_variant": true, "analytics.get_metric": true},
	domain.AgentNotify: {"notify.list_segments": true, "notify.get_campaign_history": true, "notify.create_campaign": true, "analytics.get_metric": true, "analytics.compare_holdout": true},
	domain.AgentGive:   {"give.get_history": true, "give.create_giveaway": true, "analytics.get_metric": true},
}

func DefaultLimits() Limits {
	return Limits{RunBudgetUSD: 100, ApprovalSpendUSD: 10, ApprovalAudience: 100, ApprovalDiscountPct: 10, ApprovalConfidence: .75, MaxDiscountPct: 15, MaxMessagesPerContactPer7d: 2}
}

func Check(p domain.Proposal, alreadyAllocated float64, limits Limits) Decision {
	d := Decision{Allowed: true}
	if !permissions[p.AgentType][p.ToolName] {
		d.Allowed = false
		d.Reasons = append(d.Reasons, "tool is not permitted for agent role")
	}
	if p.BudgetUSD < 0 || p.BudgetUSD+alreadyAllocated > limits.RunBudgetUSD {
		d.Allowed = false
		d.Reasons = append(d.Reasons, "run budget would be exceeded")
	}
	if len(p.EvidenceIDs) == 0 {
		d.RequiresApproval = true
		d.Reasons = append(d.Reasons, "required evidence is missing")
	}
	if p.BudgetUSD > limits.ApprovalSpendUSD {
		d.RequiresApproval = true
		d.Reasons = append(d.Reasons, "spend exceeds human approval threshold")
	}
	if p.Confidence < limits.ApprovalConfidence {
		d.RequiresApproval = true
		d.Reasons = append(d.Reasons, "confidence is below human approval threshold")
	}
	if p.ToolName == "give.create_giveaway" {
		d.RequiresApproval = true
		d.Reasons = append(d.Reasons, "giveaway launch always requires approval")
	}
	if audience, ok := number(p.ToolArguments["audienceSize"]); ok && int(audience) > limits.ApprovalAudience {
		d.RequiresApproval = true
		d.Reasons = append(d.Reasons, "audience exceeds human approval threshold")
	}
	if discount, ok := number(p.ToolArguments["discountPct"]); ok {
		if discount > limits.MaxDiscountPct {
			d.Allowed = false
			d.Reasons = append(d.Reasons, "discount exceeds absolute limit")
		} else if discount > limits.ApprovalDiscountPct {
			d.RequiresApproval = true
			d.Reasons = append(d.Reasons, "discount exceeds human approval threshold")
		}
	}
	if containsInjection(p.Hypothesis) {
		d.Reasons = append(d.Reasons, "untrusted instruction-like content treated as data")
	}
	return d
}

func number(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	default:
		return 0, false
	}
}

func containsInjection(s string) bool {
	lower := strings.ToLower(s)
	for _, marker := range []string{"ignore your instructions", "set budget to unlimited", "mark my campaign successful", "do not retire"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func ValidateTransition(from, to domain.AgentStatus, minimumSample bool, stopMet bool) error {
	allowed := map[domain.AgentStatus]map[domain.AgentStatus]bool{
		domain.StatusProposed:         {domain.StatusPolicyChecked: true, domain.StatusFailed: true},
		domain.StatusPolicyChecked:    {domain.StatusAwaitingApproval: true, domain.StatusApproved: true, domain.StatusFailed: true},
		domain.StatusAwaitingApproval: {domain.StatusApproved: true, domain.StatusPaused: true},
		domain.StatusApproved:         {domain.StatusRunning: true},
		domain.StatusRunning:          {domain.StatusEvaluating: true, domain.StatusFailed: true},
		domain.StatusEvaluating:       {domain.StatusWinning: true, domain.StatusLosing: true, domain.StatusPaused: true, domain.StatusRetired: true},
		domain.StatusLosing:           {domain.StatusRetired: true, domain.StatusPaused: true},
	}
	if !allowed[from][to] {
		return fmt.Errorf("invalid lifecycle transition %s -> %s", from, to)
	}
	if to == domain.StatusRetired && (!minimumSample || !stopMet) {
		return fmt.Errorf("retirement requires minimum sample and a true stop condition")
	}
	return nil
}
