package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("unix", "/tmp/testsock")
	if err != nil {
		fmt.Fprintln(os.Stderr, "socket connection error:", err)
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
