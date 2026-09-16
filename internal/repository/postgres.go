package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raikwar/mulelab/internal/domain"
)

type Postgres struct{ Pool *pgxpool.Pool }

func Connect(ctx context.Context, url string) (*Postgres, error) {
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	cfg.MaxConns = 5
	cfg.MinConns = 0
	cfg.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &Postgres{Pool: pool}, nil
}
func (p *Postgres) Close() {
	if p != nil && p.Pool != nil {
		p.Pool.Close()
	}
}
func (p *Postgres) Ready(ctx context.Context) error { return p.Pool.Ping(ctx) }

func (p *Postgres) Migrate(ctx context.Context, directory string) error {
	if _, err := p.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version text PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		var applied bool
		if err := p.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, name).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}
		sql, err := os.ReadFile(filepath.Join(directory, name))
		if err != nil {
			return err
		}
		tx, err := p.Pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, string(sql)); err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, name)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if err = tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) Seed(ctx context.Context) error {
	_, err := p.Pool.Exec(ctx, `INSERT INTO businesses (id,slug,name,description,scenario_seed) VALUES ('00000000-0000-0000-0000-000000000001','cosmic-cats','Cosmic Cats','A fictional creator store used for the deterministic recruiter demo.',424242) ON CONFLICT (slug) DO UPDATE SET name=EXCLUDED.name,description=EXCLUDED.description,scenario_seed=EXCLUDED.scenario_seed`)
	if err != nil {
		return err
	}
	_, err = p.Pool.Exec(ctx, `INSERT INTO business_snapshots (business_id,day,metrics) VALUES ('00000000-0000-0000-0000-000000000001',0,'{"orders":53,"conversionRate":0.017,"repeatPurchaseRate":0.09,"subscribers":420,"qualifiedSubscribers":0,"experimentSpendUsd":0}'::jsonb) ON CONFLICT (business_id,day) DO UPDATE SET metrics=EXCLUDED.metrics`)
	if err != nil {
		return err
	}
	_, err = p.Pool.Exec(ctx, `INSERT INTO model_profiles(id,supports_json_schema,supports_tools,is_free_allowed,purpose) VALUES ('openrouter/free',true,true,true,'Dynamic free route; exact resolved model is persisted'),('deterministic/demo-v1',true,true,true,'Labeled safe deterministic fallback') ON CONFLICT (id) DO NOTHING`)
	return err
}

func (p *Postgres) Business(ctx context.Context, slug string) (domain.Business, error) {
	var b domain.Business
	var metrics []byte
	err := p.Pool.QueryRow(ctx, `SELECT b.slug,b.name,b.description,b.scenario_seed,s.metrics FROM businesses b JOIN business_snapshots s ON s.business_id=b.id AND s.day=0 WHERE b.slug=$1`, slug).Scan(&b.ID, &b.Name, &b.Description, &b.Seed, &metrics)
	if err != nil {
		return b, err
	}
	if err = json.Unmarshal(metrics, &b.Metrics); err != nil {
		return b, err
	}
	return b, nil
}
