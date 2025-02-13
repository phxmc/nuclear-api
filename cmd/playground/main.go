package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/builders"
	"github.com/orewaee/nuclear-api/internal/config"
	"github.com/orewaee/nuclear-api/internal/logger"
	"github.com/orewaee/nuclear-api/internal/postgres"
	"github.com/orewaee/typedenv"
)

func main() {
	config.MustLoad()

	ctx := context.Background()

	pool := mustInitPostgresPool(ctx)

	passRepo := postgres.NewPassRepo(pool)

	log, err := logger.NewZerolog()
	if err != nil {
		panic(err)
	}

	passService := builders.NewPassServiceBuilder().
		PassRepo(passRepo).
		Log(log).
		Build()

	pass, err := passService.GetPassByAccountId(ctx, "0x89c098f0b9db90baabc835fa8a6a4a658a926f7b")
	if err != nil {
		panic(err)
	}

	fmt.Println(pass)
}

func mustInitPostgresPool(ctx context.Context) *pgxpool.Pool {
	user := typedenv.String("POSTGRES_USER")
	password := typedenv.String("POSTGRES_PASSWORD")
	addr := typedenv.String("POSTGRES_ADDR")
	database := typedenv.String("POSTGRES_DATABASE")

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, addr, database)

	pool, err := postgres.NewPool(ctx, connString)
	if err != nil {
		panic(err)
	}

	return pool
}
