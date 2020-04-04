package main

import (
	"crypto/tls"
	"log"
	"net/rpc"
)

type Result struct {
	Num, Ans int
}

type Cal int

func (cal *Cal) Square(num int, result *Result) error {
	result.Num = num
	result.Ans = num * num
	return nil
}

func main() {
	rpc.Register(new(Cal))
	cert, err := tls.LoadX509KeyPair("server/server.crt", "server/server.key")
	if err != nil {
		log.Fatal("failed to read server.cert")
	}
	// fmt.Println(cert)
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, _ := tls.Listen("tcp", ":1234", config)
	log.Printf("Serving RPC server on port %d", 1234)

	// rpc.HandleHTTP()

	// if err := http.ListenAndServe(":1234", nil); err != nil {
	// 	log.Fatal("Error serving: ", err)
	// }
	for {
		conn, _ := listener.Accept()
		defer conn.Close()
		go rpc.ServeConn(conn)
	}
}
