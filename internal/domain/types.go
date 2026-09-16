package domain

import "time"

type AgentStatus string

const (
	StatusProposed         AgentStatus = "PROPOSED"
	StatusPolicyChecked    AgentStatus = "POLICY_CHECKED"
	StatusAwaitingApproval AgentStatus = "AWAITING_APPROVAL"
	StatusApproved         AgentStatus = "APPROVED"
	StatusRunning          AgentStatus = "RUNNING"
	StatusEvaluating       AgentStatus = "EVALUATING"
	StatusWinning          AgentStatus = "WINNING"
	StatusLosing           AgentStatus = "LOSING"
	StatusPaused           AgentStatus = "PAUSED"
	StatusRetired          AgentStatus = "RETIRED"
	StatusFailed           AgentStatus = "FAILED"
)

type AgentType string

const (
	AgentSupervisor AgentType = "SUPERVISOR"
	AgentStore      AgentType = "STORE"
	AgentNotify     AgentType = "NOTIFY"
	AgentGive       AgentType = "GIVE"
)

type Metrics struct {
	Orders                 int     `json:"orders"`
	ConversionRate         float64 `json:"conversionRate"`
	RepeatPurchaseRate     float64 `json:"repeatPurchaseRate"`
	Subscribers            int     `json:"subscribers"`
	QualifiedSubscribers   int     `json:"qualifiedSubscribers"`
	ExperimentSpendUSD     float64 `json:"experimentSpendUsd"`
	UnsubscribeRate        float64 `json:"unsubscribeRate"`
	IncrementalOrders      int     `json:"incrementalOrders"`
	CostPerQualifiedSignup float64 `json:"costPerQualifiedSubscriber"`
}

type Business struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Seed        int64   `json:"seed"`
	Metrics     Metrics `json:"metrics"`
}

type Objective struct {
	ID          string  `json:"id"`
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	TargetValue float64 `json:"targetValue"`
	BudgetUSD   float64 `json:"budgetUsd"`
}

type Proposal struct {
	ID               string         `json:"id"`
	Version          int            `json:"version"`
	AgentType        AgentType      `json:"agentType"`
	Hypothesis       string         `json:"hypothesis"`
	ToolName         string         `json:"toolName"`
	ToolArguments    map[string]any `json:"toolArguments"`
	PrimaryMetric    string         `json:"primaryMetric"`
	ExpectedOutcome  string         `json:"expectedOutcome"`
	BudgetUSD        float64        `json:"budgetUsd"`
	SuccessThreshold string         `json:"successThreshold"`
	StopCondition    string         `json:"stopCondition"`
	EvidenceIDs      []string       `json:"evidenceIds"`
	Confidence       float64        `json:"confidence"`
	Status           AgentStatus    `json:"status"`
	RequiresApproval bool           `json:"requiresApproval"`
	PolicyReasons    []string       `json:"policyReasons"`
}

type Result struct {
	PrimaryMetric    string  `json:"primaryMetric"`
	Baseline         float64 `json:"baseline"`
	Observed         float64 `json:"observed"`
	SampleSize       int     `json:"sampleSize"`
	SpendUSD         float64 `json:"spendUsd"`
	Successful       bool    `json:"successful"`
	StopConditionMet bool    `json:"stopConditionMet"`
	Summary          string  `json:"summary"`
}

type Agent struct {
	ID                   string      `json:"id"`
	Type                 AgentType   `json:"type"`
	Name                 string      `json:"name"`
	Status               AgentStatus `json:"status"`
	CurrentTask          string      `json:"currentTask"`
	CumulativeImpact     string      `json:"cumulativeImpact"`
	PolicyViolationCount int         `json:"policyViolations"`
}

type Experiment struct {
	ID         string      `json:"id"`
	ProposalID string      `json:"proposalId"`
	Status     AgentStatus `json:"status"`
	ApprovedAt *time.Time  `json:"approvedAt,omitempty"`
	Result     *Result     `json:"result,omitempty"`
}

type TraceEvent struct {
	ID        string         `json:"id"`
	TraceID   string         `json:"traceId"`
	Sequence  int            `json:"sequence"`
	Category  string         `json:"category"`
	Title     string         `json:"title"`
	Summary   string         `json:"summary"`
	Evidence  map[string]any `json:"evidence"`
	CreatedAt time.Time      `json:"createdAt"`
}

func CosmicCats() Business {
	return Business{ID: "cosmic-cats", Name: "Cosmic Cats", Description: "A fictional creator store with conversion and retention opportunities.", Seed: 424242, Metrics: Metrics{Orders: 53, ConversionRate: .017, RepeatPurchaseRate: .09, Subscribers: 420}}
}
