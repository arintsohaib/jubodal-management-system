package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	up := flag.Bool("up", false, "Run migrations up")
	down := flag.Bool("down", false, "Run migrations down")
	dbURL := os.Getenv("DATABASE_URL")

	flag.Parse()

	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Ensure migrations table exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	if *up {
		runUp(ctx, pool)
	} else if *down {
		runDown(ctx, pool)
	} else {
		fmt.Println("Usage: go run cmd/migrate/main.go -up | -down")
	}
}

func runUp(ctx context.Context, pool *pgxpool.Pool) {
	files, _ := filepath.Glob("migrations/*.up.sql")
	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file)
		
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)", version).Scan(&exists)
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}

		if exists {
			continue
		}

		log.Printf("Applying migration: %s", version)
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}

		_, err = tx.Exec(ctx, string(content))
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Migration failed for %s: %v", version, err)
		}

		_, err = tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Failed to record migration: %v", err)
		}

		tx.Commit(ctx)
		log.Printf("✓ Success: %s", version)
	}
}

func runDown(ctx context.Context, pool *pgxpool.Pool) {
	var version string
	err := pool.QueryRow(ctx, "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == pgx.ErrNoRows {
		log.Println("No migrations to rollback")
		return
	}
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	downFile := strings.Replace(version, ".up.sql", ".down.sql", 1)
	filePath := filepath.Join("migrations", downFile)

	log.Printf("Rolling back: %s", version)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read rollback file: %v", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(ctx, string(content))
	if err != nil {
		tx.Rollback(ctx)
		log.Fatalf("Rollback failed for %s: %v", version, err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version=$1", version)
	if err != nil {
		tx.Rollback(ctx)
		log.Fatalf("Failed to update migration records: %v", err)
	}

	tx.Commit(ctx)
	log.Printf("✓ Rollback complete: %s", version)
}
