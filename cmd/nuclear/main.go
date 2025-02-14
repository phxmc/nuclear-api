package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/builders"
	"github.com/orewaee/nuclear-api/internal/config"
	"github.com/orewaee/nuclear-api/internal/controllers"
	"github.com/orewaee/nuclear-api/internal/disk"
	"github.com/orewaee/nuclear-api/internal/logger"
	"github.com/orewaee/nuclear-api/internal/postgres"
	"github.com/orewaee/nuclear-api/internal/redis"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/orewaee/typedenv"
	goredis "github.com/redis/go-redis/v9"
)

func main() {
	config.MustLoad()

	ctx := context.Background()

	postgresPool := mustInitPostgresPool(ctx)
	redisClient := mustInitRedisClient(ctx)

	accountRepo := postgres.NewAccountRepo(postgresPool)
	tempAccountRepo := redis.NewTempAccountRepo(redisClient)
	loginCodeRepo := redis.NewLoginCodeRepo(redisClient)
	tokenRepo := redis.NewTokenRepo(redisClient)
	staticRepo := disk.NewStaticRepo()
	passRepo := postgres.NewPassRepo(postgresPool)

	log, err := logger.NewZerolog()
	if err != nil {
		panic(err)
	}

	authApi := builders.NewAuthServiceBuilder().
		AccountRepo(accountRepo).
		LoginCodeRepo(loginCodeRepo).
		TokenRepo(tokenRepo).
		Build()

	accountApi := builders.NewAccountServiceBuilder().
		AccountRepo(accountRepo).
		TempAccountRepo(tempAccountRepo).
		Build()

	staticApi := services.NewStaticService(staticRepo)

	emailApi := services.NewEmailService(
		typedenv.String("SMTP_FROM"),
		typedenv.String("SMTP_PASSWORD"),
		typedenv.String("SMTP_HOST"),
		typedenv.String("SMTP_PORT"),
	)

	passApi := builders.NewPassServiceBuilder().
		PassRepo(passRepo).
		Log(log).
		Build()

	rest := controllers.NewRestController(typedenv.String("NUCLEAR_ADDR"), authApi, accountApi, emailApi, staticApi, passApi, log)
	if err := rest.Run(); err != nil {
		panic(err)
	}
}

func mustInitRedisClient(ctx context.Context) *goredis.Client {
	addr := typedenv.String("REDIS_ADDR")
	password := typedenv.String("REDIS_PASSWORD")

	client, err := redis.NewClient(ctx, addr, password, 0)
	if err != nil {
		panic(err)
	}

	return client
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
