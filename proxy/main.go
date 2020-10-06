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

func redirect(from string, to string) {

	remote, err := url.Parse(to)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc(from, handler(proxy))

}

func main() {

	redirect("localhost1/", "http://localhost:8080")
	redirect("localhost2/", "http://localhost:8081")

	fmt.Println("Proxy running on port 8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
