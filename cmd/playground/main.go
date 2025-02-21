package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/builders"
	"github.com/orewaee/nuclear-api/internal/config"
	"github.com/orewaee/nuclear-api/internal/logger"
	"github.com/orewaee/nuclear-api/internal/postgres"
	"github.com/orewaee/nuclear-api/internal/redis"
	"github.com/orewaee/typedenv"
	goredis "github.com/redis/go-redis/v9"
)

func main() {
	config.MustLoad()

	ctx := context.Background()

	postgresPool := mustInitPostgresPool(ctx)

	accountRepo := postgres.NewAccountRepo(postgresPool)
	passRepo := postgres.NewPassRepo(postgresPool)

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

	passApi := builders.NewPassApiBuilder().
		PassRepo(passRepo).
		Log(log).
		Build()

	/*
		email := "twd.name@gmail.com"
		tempAccount, deadline, err := accountApi.AddTempAccount(ctx, email, time.Second*30)
		if err != nil {
			panic(err)
		}

		fmt.Printf("DEADLINE: %+v\n", deadline)

		fmt.Println(accountApi.SaveTempAccount(ctx, email, tempAccount.Code))

	*/

	account, err := accountApi.GetAccountByEmail(ctx, "twd.name@gmail.com")
	if err != nil {
		panic(err)
	}

	passes, err := passApi.GetPassHistoryByAccountId(ctx, account.Id)
	if err != nil {
		panic(err)
	}

	for _, pass := range passes {
		fmt.Printf("%+v\n", pass)
	}
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
