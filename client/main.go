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
		Address:   "consul-server:8500",
		Transport: cleanhttp.DefaultPooledTransport(),
	})
	if err != nil {
		log.Fatalf("Start server failed: %v\n", err)
	}

	if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "client",
		Name:    "client",
		Address: getMyIP(),
		Port:    443,
	}); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer client.Agent().ServiceDeregister("client")

	svc, _ := connect.NewService("client", client)
	defer svc.Close()

	resp, err := svc.HTTPClient().Get("https://server.service.consul")
	if err != nil {
		log.Fatalf("Request failed: %v\n", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	msg, _ := json.Marshal(map[string]string{
		"Response": string(body),
	})
	fmt.Printf("Response: %v\n", msg)

	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatal(err)
	}
}
