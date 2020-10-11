package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"gopkg.in/yaml.v3"
)

//ProxyConfig structure contains fields related to reverse proxy configuration
type ProxyConfig struct {
	Network      string             `yaml:"network"`
	Port         string             `yaml:"port"`
	Redirections []ProxyRedirection `yaml:"redirections"`
}

//ProxyRedirection structure contains fields related to redirection
//properties and path to trust material
type ProxyRedirection struct {
	FromPattern string `yaml:"from_pattern"`
	ToURL       string `yaml:"to_url"`
	CrtPath     string `yaml:"crt_path"`
	KeyPath     string `yaml:"key_path"`
}

//ReverseProxy struct contains *http.Server and net.Listener fields
type ReverseProxy struct {
	server   *http.Server
	listener net.Listener
}

//ParseConfigFile parses YAML configuration
func ParseConfigFile(filePath string) ProxyConfig {
	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	config := ProxyConfig{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func reverseProxyHandler(rp *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		rp.ServeHTTP(w, r)
	}
}

func setUpRedirection(i int, redirection ProxyRedirection, tlsConfig *tls.Config) {
	var tlsError error

	tlsConfig.Certificates[i], tlsError = tls.LoadX509KeyPair(redirection.CrtPath, redirection.KeyPath)
	if tlsError != nil {
		log.Fatal(tlsError)
	}

	redirectedURL, err := url.Parse(redirection.ToURL)
	if err != nil {
		log.Fatal(err)
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(redirectedURL)
	http.HandleFunc(redirection.FromPattern, reverseProxyHandler(reverseProxy))
}

//NewReverseProxy create reverse proxy using TLS
func NewReverseProxy(configs ProxyConfig) *ReverseProxy {
	var err error

	reverseProxy := &ReverseProxy{}
	redirections := configs.Redirections

	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, len(redirections))

	for i := range redirections {
		setUpRedirection(i, redirections[i], tlsConfig)
	}

	tlsConfig.BuildNameToCertificate()

	reverseProxy.server = &http.Server{
		TLSConfig: tlsConfig,
	}
	reverseProxy.listener, err = tls.Listen(configs.Network, configs.Port, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	return reverseProxy
}

func main() {

	configs := ParseConfigFile("config.yaml")
	reverseProxy := NewReverseProxy(configs)
	fmt.Println("Proxy running...")
	log.Fatal(reverseProxy.server.Serve(reverseProxy.listener))
}
