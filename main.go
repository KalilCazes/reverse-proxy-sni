package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-HEADERTEST", "baby steps")
		p.ServeHTTP(w, r)
	}
}

func main() {
	remote, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))

	fmt.Println("Proxy running on port 8000...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
