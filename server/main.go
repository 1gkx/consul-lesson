package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/hashicorp/go-cleanhttp"
)

func getMyIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	return addrs[0].String()
}

func main() {

	// Create a Consul API client
	client, _ := api.NewClient(&api.Config{
		Address: "consul-server:8500",
		// Address:   "localhost:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})

	err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "server",
		Name:    "server",
		Address: getMyIP(),
		Port:    443,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer client.Agent().ServiceDeregister("server")

	// q, _, _ := client.Agent().ConnectCARoots(&api.QueryOptions{})
	// a, _, _ := client.Agent().ConnectCALeaf("server", &api.QueryOptions{})

	// cert, err := tls.X509KeyPair([]byte(a.CertPEM), []byte(a.PrivateKeyPEM))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// CA_Pool := x509.NewCertPool()
	// CA_Pool.AppendCertsFromPEM([]byte(q.Roots[0].RootCertPEM))
	// config := &tls.Config{
	// 	ClientAuth:   tls.RequireAndVerifyClientCert,
	// 	Certificates: []tls.Certificate{cert},
	// 	RootCAs:      CA_Pool,
	// 	ClientCAs:    CA_Pool,
	// 	// InsecureSkipVerify: true,
	// }

	svc, err := connect.NewService("server", client)
	if err != nil {
		log.Fatalf("Register service failed: %v\n", err)
	}
	defer svc.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
		log.Print("Certificate chain:")
		for i, cert := range r.TLS.PeerCertificates {
			subject := cert.Subject
			issuer := cert.Issuer
			log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
			log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
		}
		log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")

		msg, _ := json.Marshal(map[string]string{
			"Response": "OK",
		})
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	})
	// Creating an HTTP server that serves via Connect
	// server := &http.Server{
	// 	Addr:      ":443",
	// 	TLSConfig: config,
	// 	// ... other standard fields
	// }

	server := &http.Server{
		Addr: ":443",
		// Handler:   apiHandler{},
		TLSConfig: svc.ServerTLSConfig(),
	}

	// Serve!
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Start server failed: %v\n", err)
	}
}
