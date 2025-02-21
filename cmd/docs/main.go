package main

import (
	"fmt"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger", func(writer http.ResponseWriter, request *http.Request) {
		content, err := os.ReadFile("./api/swagger.json")
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(writer, string(content))
	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "http://127.0.0.1:8040/swagger",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Nuclear Docs",
			},
			Theme:    scalar.ThemeDeepSpace,
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

	http.ListenAndServe(":8040", mux)
}
