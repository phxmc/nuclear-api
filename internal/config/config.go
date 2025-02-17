package config

import (
	"github.com/joho/godotenv"
	"github.com/orewaee/typedenv"
	"os"
	"time"
)

func MustLoad() {
	envs := []string{
		"config/.env",
		"config/postgres.env",
		"config/redis.env",
	}

	err := godotenv.Load(envs...)

	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	typedenv.DefaultString("NUCLEAR_ADDR", ":8080")

	typedenv.DefaultDuration("ACCESS_LIFETIME", time.Minute*10)
	typedenv.DefaultDuration("REFRESH_LIFETIME", time.Hour*24)

	typedenv.DefaultString("ALPHABET", "abcdefghijklmnopqrstuvwxyz")

	typedenv.DefaultDuration("TEMP_ACCOUNT_LIFETIME", time.Minute*5)
}
