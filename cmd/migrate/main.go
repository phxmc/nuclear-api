package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/orewaee/typedenv"
	"os"
)

func getUrl() string {
	user := typedenv.String("POSTGRES_USER")
	password := typedenv.String("POSTGRES_PASSWORD")
	host := typedenv.String("POSTGRES_HOST")
	port := typedenv.String("POSTGRES_PORT")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/nuclear?sslmode=disable",
		user, password, host, port)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("invalid syntax")
		return
	}

	action := os.Args[1]

	err := godotenv.Load("config/postgres.env")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	url := getUrl()

	m, err := migrate.New("file://migrations", url)
	if err != nil {
		panic(err)
	}

	var f func() error

	switch action {
	case "up":
		f = m.Up
		break
	case "down":
		f = m.Down
	default:
		fmt.Println("invalid syntax")
		return
	}

	if err := f(); err != nil {
		panic(err)
	}

	fmt.Println("complete")
}
