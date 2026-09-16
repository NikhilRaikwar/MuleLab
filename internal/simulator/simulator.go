package simulator

import (
	"fmt"
	"math/rand"

	"github.com/raikwar/mulelab/internal/domain"
)

type Engine struct{ Seed int64 }

func (e Engine) Evaluate(p domain.Proposal, days int) domain.Result {
	if days < 1 {
		days = 1
	}
	rng := rand.New(rand.NewSource(e.Seed + hash(p.ID)))
	switch p.AgentType {
	case domain.AgentStore:
		lift := .008 + float64(rng.Intn(2))/1000
		return domain.Result{PrimaryMetric: "store_conversion_rate", Baseline: .017, Observed: .017 + lift, SampleSize: 720 * days, SpendUSD: p.BudgetUSD, Successful: lift >= .004, Summary: fmt.Sprintf("Conversion improved from 1.7%% to %.1f%%.", (.017+lift)*100)}
	case domain.AgentNotify:
		orders := 7 + rng.Intn(3)
		return domain.Result{PrimaryMetric: "incremental_orders_vs_holdout", Baseline: 0, Observed: float64(orders), SampleSize: 180, SpendUSD: p.BudgetUSD, Successful: orders >= 5, Summary: fmt.Sprintf("Generated %d incremental orders versus holdout.", orders)}
	default:
		qualified := 5 + rng.Intn(2)
		cost := p.BudgetUSD / float64(qualified)
		return domain.Result{PrimaryMetric: "cost_per_qualified_subscriber", Baseline: 2.1, Observed: cost, SampleSize: 160, SpendUSD: p.BudgetUSD, Successful: cost <= 1.5, StopConditionMet: cost >= 3, Summary: fmt.Sprintf("Qualified subscriber cost was $%.2f; deterministic stop threshold is $3.00.", cost)}
	}
}

func hash(s string) int64 {
	var h int64
	for _, r := range s {
		h = h*31 + int64(r)
	}
	return h
}
