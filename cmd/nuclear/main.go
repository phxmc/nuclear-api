package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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

	log, err := logger.NewZerolog()
	if err != nil {
		panic(err)
	}

	authApi := services.NewAuthServiceBuilder().
		AccountRepo(accountRepo).
		LoginCodeRepo(loginCodeRepo).
		TokenRepo(tokenRepo).
		Build()

	accountApi := services.NewAccountServiceBuilder().
		AccountRepo(accountRepo).
		TempAccountRepo(tempAccountRepo).
		Build()

	staticApi := services.NewStaticService(staticRepo)

	emailApi := services.NewEmailService("orewaee@gmail.com", typedenv.String("SMTP_PASSWORD"), "smtp.gmail.com", "587")

	rest := controllers.NewRestController(":8080", authApi, accountApi, emailApi, staticApi, log)
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
	pool, err := postgres.NewPool(ctx, "postgres://root:kYJuSfL4FX7Mtcy2badzHpn9GmqUve6r@localhost:5442/nuclear?sslmode=disable")
	if err != nil {
		panic(err)
	}

	return pool
}
