package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, h.testHandler)
}

type Handler struct {
	testHandler string
}

func Newhandler(h string) *Handler {
	newHandler := &Handler{
		testHandler: h,
	}
	return newHandler
}

func TestRedirection(t *testing.T) {
	//setup reverse proxy
	reverseProxy := NewReverseProxy(ParseConfigFile("../config-test.yaml"))
	proxyServer := &httptest.Server{
		Listener: reverseProxy.listener,
		Config:   reverseProxy.server,
	}
	proxyServer.Start()
	defer proxyServer.Close()

	config := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: config}

	//setup fake backend server
	fakeConfig := ParseConfigFile("../config-test.yaml")
	listener, err := net.Listen(fakeConfig.Network, fakeConfig.Redirections[0].ToURL[7:])
	if err != nil {
		log.Fatal(err)
	}
	backendServer := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: Newhandler("Foward")},
	}

	backendServer.Start()
	defer proxyServer.Close()

	r, err := client.Get("https://localhost1:25000")
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, bodyString, "Foward", "The body should be 'Foward'")
	assert.Equal(t, r.TLS.PeerCertificates[0].Issuer.CommonName, "localhost1", "The certificate's common name should be 'localhost1'")
	assert.Equal(t, r.StatusCode, http.StatusOK, "The HTTP status should be 200")
}

func TestConfigParser(t *testing.T) {
	config := ParseConfigFile("../config-test.yaml")

	assert.Equal(t, config.Network, "tcp")
	assert.Equal(t, config.Port, ":25000")

	assert.Equal(t, config.Redirections[0].ToURL, "http://localhost:9999")
	assert.Equal(t, config.Redirections[0].FromPattern, "localhost1/")
	assert.Equal(t, config.Redirections[0].CrtPath, "../proxy/localhost1.crt")
	assert.Equal(t, config.Redirections[0].KeyPath, "../proxy/localhost1.key")

	assert.Equal(t, config.Redirections[1].ToURL, "http://localhost:1111")
	assert.Equal(t, config.Redirections[1].FromPattern, "localhost2/")
	assert.Equal(t, config.Redirections[1].CrtPath, "../proxy/localhost2.crt")
	assert.Equal(t, config.Redirections[1].KeyPath, "../proxy/localhost2.key")
}

func TestNewReverseProxy(t *testing.T) {

	http.DefaultServeMux = new(http.ServeMux)

	reverseProxy := NewReverseProxy(ParseConfigFile("../config-test.yaml"))

	assert.Equal(t, reverseProxy.listener.Addr().Network(), "tcp")
	assert.Equal(t, len(reverseProxy.server.TLSConfig.Certificates), 2)
	assert.True(t, strings.Contains(reverseProxy.listener.Addr().String(), "25000"))

}
