package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloServer)
	fmt.Println("Application 2 running on port 8081...")
	http.ListenAndServe(":8081", nil)
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "App2")
	log.Println(r.URL)
}
