package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/raikwar/mulelab/internal/domain"
	"github.com/raikwar/mulelab/internal/policy"
)

func TestApprovalAndIdempotency(t *testing.T) {
	r := New(policy.DefaultLimits())
	p := domain.Proposal{AgentType: domain.AgentGive, ToolName: "give.create_giveaway", BudgetUSD: 25, EvidenceIDs: []string{"give:history"}, Confidence: .88}
	call := Call{Name: p.ToolName, IdempotencyKey: "run-1:give-1:v1", Proposal: p}
	if _, err := r.Execute(context.Background(), call, false); err == nil {
		t.Fatal("giveaway ran without approval")
	}
	if _, err := r.Execute(context.Background(), call, true); err != nil {
		t.Fatalf("approved call failed: %v", err)
	}
	if _, err := r.Execute(context.Background(), call, true); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate was not detected: %v", err)
	}
}
