package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const port = ":8000"
const proto = "tcp"

//ReverseProxy struct contains *http.Server and net.Listener fields
type ReverseProxy struct {
	server   *http.Server
	listener net.Listener
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)
		p.ServeHTTP(w, r)
	}
}

func redirect(from string, to string) {

	redirectedURL, err := url.Parse(to)
	if err != nil {
		log.Fatal(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(redirectedURL)
	http.HandleFunc(from, handler(reverseProxy))

}

func createCertificate(cert string, key string) tls.Certificate {

	appCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	return appCert

}

func (rp *ReverseProxy) setUpTLS() {
	var err error
	var cert = []string{"localhost1.crt", "localhost2.crt"}
	var key = []string{"localhost1.key", "localhost2.key"}

	appCert1 := createCertificate(cert[0], key[0])
	appCert2 := createCertificate(cert[1], key[1])

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{appCert1, appCert2}}

	tlsConfig.BuildNameToCertificate()

	rp.server = &http.Server{
		TLSConfig: tlsConfig,
	}

	rp.listener, err = tls.Listen(proto, port, tlsConfig)

	if err != nil {
		log.Fatal(err)
	}

}

//NewReverseProxy create reverse proxy using TLS
func NewReverseProxy() *ReverseProxy {

	reverseProxy := &ReverseProxy{}

	redirect("localhost1/", "http://localhost:8080")
	redirect("localhost2/", "http://localhost:8081")

	reverseProxy.setUpTLS()

	return reverseProxy
}

func main() {

	reverseProxy := NewReverseProxy()
	fmt.Println("Proxy running on port 8000...")
	log.Fatal(reverseProxy.server.Serve(reverseProxy.listener))
}
