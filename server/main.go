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

	port := 3010

	// Create a Consul API client
	client, err := api.NewClient(&api.Config{
		Address:   "demo-server-agent:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		CheckID:  "0000-0001",
		Name:     "CheckName",
		Interval: "10s",
		TLSSkipVerify: true,
		TCP:           fmt.Sprintf("%s:%d", getMyIP(), port),
	}

	if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "server",
		Name:    "server",
		Address: getMyIP(),
		Port:    port,
		Check:   check,
		Connect: &api.AgentServiceConnect{
			Native: true,
		},
	}); err != nil {
		log.Fatalf("Register service failed: %v\n", err)
	}
	defer unreg(client.Agent())

	svc, err := connect.NewService("server", client)
	if err != nil {
		log.Fatalf("Get TLS config failed: %v\n", err)
	}
	defer svc.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request: %+v\n", r)
		msg, _ := json.Marshal(map[string]string{
			"Response": "OK",
		})
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	})

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: svc.ServerTLSConfig(),
	}

	// Serve!
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Start server failed: %v\n", err)
	}
}

func unreg(agent *api.Agent) {
	if err := agent.CheckDeregister("0000-0001"); err != nil {
		fmt.Printf("Check deregister: %v\n", err)
	}
	if err := agent.ServiceDeregister("server"); err != nil {
		fmt.Printf("Service deregister: %v\n", err)
	}
}
