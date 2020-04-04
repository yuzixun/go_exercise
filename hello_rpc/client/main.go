package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/rpc"
)

type Result struct {
	Num, Ans int
}

func main() {
	// cert, err := tls.LoadX509KeyPair("client/client.crt", "client/client.key")
	if err != nil {
		log.Fatal("failed to load x509 key pair")
	}

	certPool := x509.NewCertPool()
	certBytes, err := ioutil.ReadFile("server/server.crt")
	if err != nil {
		log.Fatal("failed to read server.cert")
	}

	certPool.AppendCertsFromPEM(certBytes)

	config := &tls.Config{
		// Certificates: []tls.Certificate{cert},
		RootCAs: certPool,
	}

	// config := &tls.Config{
	// 	InsecureSkipVerify: true,
	// }
	conn, err := tls.Dial("tcp", "localhost:1234", config)
	defer conn.Close()
	client := rpc.NewClient(conn)

	var result Result
	// client, _ := rpc.DialHTTP("tcp", "localhost:1234")
	if err := client.Call("Cal.Square", 12, &result); err != nil {
		log.Fatal("failed to call Cal.Square. ", err)
	}

	// asyncCall := client.Go("Cal.Square", 12, &result, nil)
	// log.Printf("%d^2 = %d", result.Num, result.Ans)
	// <-asyncCall.Done
	log.Printf("%d^2 = %d", result.Num, result.Ans)
}
