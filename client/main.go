package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	client, err := api.NewClient(&api.Config{
		Address: "localhost:8500",
		// Address:   "consul-server:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})
	if err != nil {
		log.Fatalf("Start server failed: %v\n", err)
	}

	check := &api.AgentServiceCheck{
		CheckID:       "check client",
		Name:          "CheckName",
		Interval:      "3s",
		Timeout:       "15s",
		TLSServerName: "client.service.consul",
		TCP:           fmt.Sprintf("%s:%s", getMyIP(), "443"),
	}

	if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "client",
		Name:    "client",
		Address: getMyIP(),
		Port:    443,
		Check:   check,
	}); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer client.Agent().ServiceDeregister("client")

	svc, _ := connect.NewService("client", client)
	defer svc.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.HTTPClient().Get("https://server.service.consul")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		msg, _ := json.Marshal(map[string]string{
			"Response": string(body),
		})

		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	})

	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatal(err)
	}
}
