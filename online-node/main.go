package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

func getMyIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	return addrs[0].String()
}

func main() {

	client, err := api.NewClient(&api.Config{
		// Address: "consul-server:8500",
		Address:   "localhost:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})
	if err != nil {
		log.Fatalf("Start server failed: %v\n", err)
	}

	if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "service-1",
		Name:    "service-1",
		Address: getMyIP(),
		Port:    3010,
	}); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer client.Agent().ServiceDeregister("service-1")

	q, _, _ := client.Agent().ConnectCARoots(&api.QueryOptions{})
	CA_Pool := x509.NewCertPool()
	CA_Pool.AppendCertsFromPEM([]byte(q.Roots[0].RootCertPEM))
	config := &tls.Config{
		RootCAs: CA_Pool,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: config,
			},
		}

		resp, err := client.Get("https://service-2.service.consul:3011")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		msg, err := json.Marshal(map[string]string{
			"Response": string(body),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	})
	if err := http.ListenAndServe(":3010", nil); err != nil {
		log.Fatal(err)
	}
}
