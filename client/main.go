package main

import (
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

	client, err := api.NewClient(&api.Config{
		Address: "consul-server:8500",
		// Address:   "localhost:8500",
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

	se, q, e := client.Health().Connect("client", "test", true, &api.QueryOptions{})
	fmt.Printf("HealthCheck: %v\nQueryMeta: %+v\nerror: %v\n", se, q, e)

	hc, q, e := client.Health().Checks("client", &api.QueryOptions{})
	fmt.Printf("HealthCheck: %v\nQueryMeta: %+v\nerror: %v\n", hc, q, e)

	// q, _, _ := client.Agent().ConnectCARoots(&api.QueryOptions{})
	// a, _, _ := client.Agent().ConnectCALeaf("client", &api.QueryOptions{})

	// Load client cert
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

	svc, _ := connect.NewService("client", client)
	defer svc.Close()

	// resp, err := svc.HTTPClient().Get("https://server.service.consul")
	// if err != nil {
	// 	log.Fatalf("Request failed: %v\n", err)
	// }
	// log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	// log.Print("Certificate chain:")
	// for i, cert := range resp.TLS.PeerCertificates {
	// 	subject := cert.Subject
	// 	issuer := cert.Issuer
	// 	log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
	// 	log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	// }
	// log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")

	// body, _ := ioutil.ReadAll(resp.Body)
	// msg, _ := json.Marshal(map[string]string{
	// 	"Response": string(body),
	// })
	// fmt.Printf("Response: %v\n", msg)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	// 	client := &http.Client{
	// 		Transport: &http.Transport{
	// 			TLSClientConfig: config,
	// 		},
	// 	}

	// 	resp, err := client.Get("https://server.service.consul")
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}

	// 	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	// 	log.Print("Certificate chain:")
	// 	for i, cert := range resp.PeerCertificates {
	// 		subject := cert.Subject
	// 		issuer := cert.Issuer
	// 		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
	// 		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	// 	}
	// 	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}

	// 	msg, err := json.Marshal(map[string]string{
	// 		"Response": string(body),
	// 	})
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write(msg)
	// })
	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatal(err)
	}
}
