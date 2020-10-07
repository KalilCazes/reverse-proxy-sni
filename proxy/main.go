package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)
		p.ServeHTTP(w, r)
	}
}

func redirect(from string, to string) {

	remote, err := url.Parse(to)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc(from, handler(proxy))

}

func main() {
	var err error

	appCert1, err := tls.LoadX509KeyPair("localhost1.crt", "localhost1.key")
	if err != nil {
		log.Fatal(err)
	}

	appCert2, err := tls.LoadX509KeyPair("localhost2.crt", "localhost2.key")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{appCert1, appCert2}}

	tlsConfig.BuildNameToCertificate()

	proxy := &http.Server{
		TLSConfig: tlsConfig,
	}

	redirect("localhost1/", "http://localhost:8080")
	redirect("localhost2/", "http://localhost:8081")

	listener, err := tls.Listen("tcp", ":8000", tlsConfig)

	fmt.Println("Proxy running on port 8000...")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(proxy.Serve(listener))
}
