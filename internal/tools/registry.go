package tools

import (
	"context"
	"errors"
	"sync"

	"github.com/raikwar/mulelab/internal/domain"
	"github.com/raikwar/mulelab/internal/policy"
)

var ErrDuplicate = errors.New("tool call already executed for idempotency key")

type Call struct {
	Name, IdempotencyKey string
	Proposal             domain.Proposal
}
type Output struct {
	Applied bool
	Summary string
}

type Registry struct {
	mu        sync.Mutex
	completed map[string]Output
	limits    policy.Limits
}

func New(limits policy.Limits) *Registry {
	return &Registry{completed: map[string]Output{}, limits: limits}
}

func (r *Registry) Execute(ctx context.Context, call Call, approved bool) (Output, error) {
	select {
	case <-ctx.Done():
		return Output{}, ctx.Err()
	default:
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if prior, exists := r.completed[call.IdempotencyKey]; exists {
		return prior, ErrDuplicate
	}
	decision := policy.Check(call.Proposal, 0, r.limits)
	if !decision.Allowed {
		return Output{}, errors.New("deterministic policy denied tool call")
	}
	if decision.RequiresApproval && !approved {
		return Output{}, errors.New("human approval required")
	}
	result := Output{Applied: true, Summary: "simulated action accepted through typed tool boundary"}
	r.completed[call.IdempotencyKey] = result
	return result, nil
}
