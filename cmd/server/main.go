package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/lavish-gambhir/dashbeam/internal/config"
)

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "===dashbeam===")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to laod config: %v", err)
	}

	http.HandleFunc("/", index)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
