# Security model

## Trust boundaries

All product descriptions, customer text, campaign copy, retrieved notes, model output, and GraphQL input are untrusted. They are data, never instructions that expand authority.

The LLM may select among declared agent workstreams and propose typed experiments. It may not approve actions, change budgets, create tools, fabricate simulator metrics, mutate protected state, edit eval truth, or decide final retirement.

Deterministic code owns schema validation, permissions, budgets, frequency limits, idempotency, approval requirements, metric computation, state transitions, and release gates.

## Controls

- Server-only model and database credentials.
- Exact tool allowlists per agent role.
- JSON Schema validation with `additionalProperties: false`.
- Explicit limits on strings, audience, discounts, spend, calls, iterations, and run duration.
- Approval records bind the proposal version and exact arguments.
- Append-only trace events with secret redaction.
- Idempotency keys on side-effecting calls.
- No arbitrary URLs, shell execution, code execution, or private Sticker Mule API access.
- Prompt-injection golden cases must pass before release.

## Known limitations

MuleLab v0.1 is a fictional public demo, not production multi-tenancy. Rate limiting and GraphQL depth/complexity controls reduce abuse but do not replace a full identity and tenancy design.

