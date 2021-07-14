package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

func main() {

	client, _ := api.NewClient(&api.Config{
		// Address: "consul-server:8500",
		Address:   "localhost:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})

	err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "service-1",
		Name:    "service-1",
		Address: "service-1:3010",
	})
	if err != nil {
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

		resp, err := client.Get("http://service-2.service.consul")
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
