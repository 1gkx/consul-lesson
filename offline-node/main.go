package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

func LoadTlsCredentialsFromBytes(
	rootCertificate []byte,
	ownCertificate []byte,
	privateKey []byte,
) (*tls.Config, error) {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(rootCertificate)

	keyPair, err := tls.X509KeyPair(ownCertificate, privateKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{keyPair},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		RootCAs:      certPool,
	}

	return tlsConfig, nil
}

func getMyIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	return addrs[0].String()
}

func main() {

	// Create a Consul API client
	client, _ := api.NewClient(&api.Config{
		// Address: "consul-server:8500",
		Address:   "localhost:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})

	err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "service-2",
		Name:    "service-2",
		Address: getMyIP(),
		Port:    3011,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer client.Agent().ServiceDeregister("service-2")

	q, _, _ := client.Agent().ConnectCARoots(&api.QueryOptions{})
	a, _, _ := client.Agent().ConnectCALeaf("service-2", &api.QueryOptions{})

	tls, err := LoadTlsCredentialsFromBytes(
		[]byte(q.Roots[0].RootCertPEM),
		[]byte(a.CertPEM),
		[]byte(a.PrivateKeyPEM),
	)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(map[string]string{
			"Response": "OK",
		})
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	})
	// Creating an HTTP server that serves via Connect
	server := &http.Server{
		Addr:      ":3011",
		TLSConfig: tls,
		// ... other standard fields
	}

	// Serve!
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Start server failed: %v\n", err)
	}
}
