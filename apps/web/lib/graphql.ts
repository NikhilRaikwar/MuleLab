const endpoint = process.env.NEXT_PUBLIC_GRAPHQL_URL ?? "http://localhost:8080/graphql";

export async function graphql<T>(query: string, variables?: Record<string, unknown>): Promise<T> {
  const response = await fetch(endpoint, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ query, variables }),
    cache: "no-store",
  });
  const body = await response.json() as { data?: T; errors?: { message: string }[] };
  if (!response.ok || body.errors?.length || !body.data) throw new Error(body.errors?.[0]?.message ?? "GraphQL request failed");
  return body.data;
}

export type Proposal = { id: string; agentType: string; hypothesis: string; toolName: string; toolArgumentsJson: string; primaryMetric: string; expectedOutcome: string; budgetUsd: number; successThreshold: string; stopCondition: string; evidenceIds: string[]; confidence: number; status: string; requiresApproval: boolean; policyReasons: string[] };
export type Agent = { id: string; type: string; name: string; status: string; currentTask: string; cumulativeImpact: string };
export type Experiment = { id: string; proposal: Proposal; status: string; result?: { primaryMetric: string; observed: number; summary: string; successful: boolean; stopConditionMet: boolean } };
export type Run = { id: string; traceId: string; status: string; day: number; agents: Agent[]; proposals: Proposal[]; experiments: Experiment[]; modelStats: { requestedModel: string; resolvedModel: string; source: string; latencyMs: number; promptTokens: number; completionTokens: number; estimatedCostUsd: number; retries: number }[] };
