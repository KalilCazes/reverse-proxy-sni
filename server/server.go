package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloServer)
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "testing, %s!\n", r.URL.Path[1:])
	log.Println(r.URL)
}
