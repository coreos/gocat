package main

import (
	"crypto/tls"
	"fmt"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "private.key")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cert load failed:", err)
	}

	cfg := &tls.Config {
		Certificates: []tls.Certificate{cert},
		ServerName: "test_server",
		/* ClientAuth: tls.RequireAndVerifyClientCert, */
		ClientAuth: tls.RequireAnyClientCert,
	}
	cfg.BuildNameToCertificate()

	ln, err := tls.Listen("tcp", ":8080", cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen failed:", err)
		os.Exit(-1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "accept failure:", err)
			continue
		}

		fmt.Fprintln(conn, "success!")
		conn.Close()
	}

}
