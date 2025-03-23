package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/builders"
	"github.com/orewaee/nuclear-api/internal/config"
	"github.com/orewaee/nuclear-api/internal/logger"
	"github.com/orewaee/nuclear-api/internal/postgres"
	"github.com/orewaee/nuclear-api/internal/redis"
	"github.com/orewaee/typedenv"
	goredis "github.com/redis/go-redis/v9"
	"time"
)

func main() {
	config.MustLoad()

	ctx := context.TODO()

	postgresPool := mustInitPostgresPool(ctx)

	accountRepo := postgres.NewAccountRepo(postgresPool)

	redisClient := mustInitRedisClient(ctx)

	tempAccountRepo := redis.NewTempAccountRepo(redisClient)

	log, err := logger.NewZerolog()
	if err != nil {
		panic(err)
	}

	accountApi := builders.NewAccountApiBuilder().
		AccountRepo(accountRepo).
		TempAccountRepo(tempAccountRepo).
		Log(log).
		Build()

	_, _ = accountApi.GetAccountById(ctx, "")

	telegramRepo := redis.NewTelegramRepo(redisClient)

	fmt.Println(telegramRepo.SetChatState(ctx, 123, domain.StateEnterCode, time.Second*15))
	fmt.Println(telegramRepo.GetChatState(ctx, 123))
}

func mustInitPostgresPool(ctx context.Context) *pgxpool.Pool {
	user := typedenv.String("POSTGRES_USER")
	password := typedenv.String("POSTGRES_PASSWORD")
	host := typedenv.String("POSTGRES_HOST")
	port := typedenv.String("POSTGRES_PORT")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/nuclear?sslmode=disable", user, password, host, port)

	pool, err := postgres.NewPool(ctx, connString)
	if err != nil {
		panic(err)
	}

	return pool
}

func mustInitRedisClient(ctx context.Context) *goredis.Client {
	host := typedenv.String("REDIS_HOST")
	port := typedenv.String("REDIS_PORT")
	password := typedenv.String("REDIS_PASSWORD")

	addr := fmt.Sprintf("%s:%s", host, port)

	client, err := redis.NewClient(ctx, addr, password, 0)
	if err != nil {
		panic(err)
	}

	return client
}
