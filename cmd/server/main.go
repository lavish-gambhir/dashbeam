package main

import (
	"fmt"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "===dashbeam===")
}

func main() {
	http.HandleFunc("/", index)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
