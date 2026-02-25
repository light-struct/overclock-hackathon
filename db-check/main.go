package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// defaultDSN — URL для Neon, который вы указали.
const defaultDSN = "postgresql://neondb_owner:npg_AERK0Ho4fgws@ep-cold-surf-alqoktoy-pooler.c-3.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Берём DATABASE_URL из окружения, если он не задан — используем Neon DSN по умолчанию.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDSN
	}

	log.Println("Connecting to Postgres...")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(ctx)

	// 1. Проверка версии
	var version string
	if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query version() failed: %v", err)
	}
	log.Println("Connected to:", version)

	// 2. Проверка таблицы users
	log.Println("Reading users...")
	rows, err := conn.Query(ctx, `
        SELECT id, email, username, role, created_at
        FROM users
        ORDER BY id
        LIMIT 10
    `)
	if err != nil {
		log.Fatalf("Query users failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int64
			email     string
			username  string
			role      string
			createdAt time.Time
		)
		if err := rows.Scan(&id, &email, &username, &role, &createdAt); err != nil {
			log.Fatalf("Scan user row failed: %v", err)
		}
		fmt.Printf("User #%d: %s (%s), role=%s, created_at=%s\n",
			id, username, email, role, createdAt.Format(time.RFC3339))
	}
	if rows.Err() != nil {
		log.Fatalf("Rows error: %v", rows.Err())
	}

	// 3. Проверка количества попыток тестов
	var attemptsCount int64
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM test_attempts").Scan(&attemptsCount); err != nil {
		log.Fatalf("Query test_attempts count failed: %v", err)
	}
	log.Printf("Total test_attempts: %d\n", attemptsCount)

	log.Println("DB check finished successfully.")
}

