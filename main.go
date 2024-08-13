package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := git config --global user.name "${GITHUB_USERNAME}"http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
