package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("unix", "/tmp/testsock")
	defer ln.Close()

	if err != nil {
		fmt.Fprintln(os.Stderr, "socket listen error:", err)
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Fprintln(os.Stderr, "accept failure:", err)
		return
	}
	ioConnect(conn)
}

func ioConnect(rw io.ReadWriteCloser) {
	go func() {
		io.Copy(os.Stdin, rw)
		fmt.Fprintln(os.Stderr, "(input closed)")
	}()
	io.Copy(rw, os.Stdout)
	fmt.Fprintln(os.Stderr, "(closing output)")
	rw.Close()
}
