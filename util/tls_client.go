package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"io"
)

func main() {
	cert, err := tls.LoadX509KeyPair("client_cert.pem", "client_private.key")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cert load failed:", err)
		return
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "localhost:8080", cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial failed:", err)
		return
	}

	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}
