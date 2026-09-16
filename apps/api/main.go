package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/raikwar/mulelab/apps/api/graph"
	"github.com/raikwar/mulelab/internal/config"
	"github.com/raikwar/mulelab/internal/repository"
	"github.com/raikwar/mulelab/internal/runtime"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := repository.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer db.Close()
	if err := db.Migrate(ctx, "migrations"); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	if err := db.Seed(ctx); err != nil {
		log.Fatalf("database seed failed: %v", err)
	}
	rt := runtime.New(db, cfg)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Runtime: rt}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](100))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	mux := http.NewServeMux()
	mux.Handle("/graphql", withCORS(srv))
	mux.Handle("/", playground.Handler("MuleLab GraphQL", "/graphql"))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := rt.Ready(req.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready","database":"unavailable"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"ready","database":"ok"}`))
	})

	server := &http.Server{Addr: ":" + cfg.Port, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 95 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("mulelab api listening on :%s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}

func withCORS(next http.Handler) http.Handler {
	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		allowed = "http://localhost:3000"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
