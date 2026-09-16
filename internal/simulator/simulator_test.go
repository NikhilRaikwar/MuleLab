package simulator

import (
	"reflect"
	"testing"

	"github.com/raikwar/mulelab/internal/domain"
)

func TestReproducible(t *testing.T) {
	p := domain.Proposal{ID: "give-1", AgentType: domain.AgentGive, BudgetUSD: 25}
	a, b := (Engine{Seed: 424242}).Evaluate(p, 7), (Engine{Seed: 424242}).Evaluate(p, 7)
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("same seed produced different outcomes: %#v %#v", a, b)
	}
}

func TestGiveAgentCanHitDeterministicStop(t *testing.T) {
	p := domain.Proposal{ID: "give-1", AgentType: domain.AgentGive, BudgetUSD: 25}
	r := (Engine{Seed: 424242}).Evaluate(p, 7)
	if !r.StopConditionMet || r.Successful {
		t.Fatalf("expected losing giveaway, got %#v", r)
	}
}
