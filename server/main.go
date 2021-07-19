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
		Address:   "localhost:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})

	check := &api.AgentServiceCheck{
		CheckID:       "0000-0001",
		Name:          "CheckName",
		Interval:      "10s",
		TLSSkipVerify: true,
		TLSServerName: "server.service.consul",
		TCP:           fmt.Sprintf("%s:%s", getMyIP(), "443"),
	}

	err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "server",
		Name:    "server",
		Address: getMyIP(),
		Port:    443,
		Check:   check,
		Connect: &api.AgentServiceConnect{
			Native: true,
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer unreg(client.Agent())

	svc, err := connect.NewService("server", client)
	if err != nil {
		log.Fatalf("Register service failed: %v\n", err)
	}
	defer svc.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(map[string]string{
			"Response": "OK",
		})
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	})

	tls := svc.ServerTLSConfig()
	tls.InsecureSkipVerify = true
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tls,
	}

	// Serve!
	if err := server.ListenAndServeTLS("", ""); err != nil {
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
