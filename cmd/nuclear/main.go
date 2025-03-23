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
	"github.com/orewaee/nuclear-api/internal/telegram"
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
	nicknameRepo := postgres.NewNicknameRepo(postgresPool)
	telegramRepo := redis.NewTelegramRepo(redisClient)

	log, err := logger.NewZerolog()
	if err != nil {
		panic(err)
	}

	authApi := builders.NewAuthServiceBuilder().
		AccountRepo(accountRepo).
		LoginCodeRepo(loginCodeRepo).
		TokenRepo(tokenRepo).
		Log(log).
		Build()

	accountApi := builders.NewAccountApiBuilder().
		AccountRepo(accountRepo).
		TempAccountRepo(tempAccountRepo).
		Log(log).
		Build()

	passApi := builders.NewPassApiBuilder().
		PassRepo(passRepo).
		Log(log).
		Build()

	nicknameApi := builders.NewNicknameApiBuilder().
		NicknameRepo(nicknameRepo).
		Log(log).
		Build()

	telegramApi := builders.NewTelegramApiBuilder().
		TelegramRepo(telegramRepo).
		Log(log).
		Build()

	staticApi := services.NewStaticService(staticRepo)

	emailApi := services.NewEmailService(
		typedenv.String("SMTP_FROM"),
		typedenv.String("SMTP_PASSWORD"),
		typedenv.String("SMTP_HOST"),
		typedenv.String("SMTP_PORT"),
	)

	bot := telegram.NewBot(accountApi, telegramApi, emailApi, log)

	go func() {
		err = bot.Run(ctx, typedenv.String("TELEGRAM_TOKEN"))
		if err != nil {
			panic(err)
		}
	}()

	rest := controllers.NewRestController(typedenv.String("NUCLEAR_ADDR"), authApi, accountApi, emailApi, staticApi, passApi, nicknameApi, log)
	if err := rest.Run(); err != nil {
		panic(err)
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
