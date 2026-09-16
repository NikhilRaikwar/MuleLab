INSERT INTO businesses (id, slug, name, description, scenario_seed)
VALUES ('00000000-0000-0000-0000-000000000001', 'cosmic-cats', 'Cosmic Cats', 'A fictional creator store used for the deterministic recruiter demo.', 424242);

INSERT INTO business_snapshots (business_id, day, metrics)
VALUES ('00000000-0000-0000-0000-000000000001', 0, '{"products":12,"visitors30d":3100,"conversionRate":0.017,"orders30d":53,"repeatPurchaseRate":0.09,"subscribers":420,"customers":86}');

INSERT INTO model_profiles (id, supports_json_schema, supports_tools, is_free_allowed, purpose)
VALUES ('openrouter/free', true, true, true, 'Dynamic free route; exact resolved model is recorded'),
       ('deterministic/demo-v1', true, true, true, 'Labeled recruiter-safe fallback and CI fixture');

