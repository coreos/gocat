package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"syscall"
)

//TODO: review error handling

func flagAssert(pred bool, display string) {
	if !pred {
		fmt.Fprintln(os.Stderr, display)
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	key := flag.String("key", "", "path to private key")
	cert := flag.String("cert", "", "path to cert pem")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "gocat [unix socket file] [host]:port")
		fmt.Fprintln(os.Stderr, "  (when gocat is socket activated, the second argument is ignored)")
		flag.PrintDefaults()
	}
	flag.Parse()
	flagAssert(flag.NArg() == 2, "error: gocat requires exactly 2 arguments")
	flagAssert(*key == "" && *cert == "" || *key != "" && *cert != "",
		"error: gocat requires both a key and a certificate to be specified")
	socket := flag.Arg(0)
	server := flag.Arg(1)

	fmt.Println(socket, server, *cert, *key) //TODO: remove
	gocat(socket, server, *cert, *key)
}

func gocat(socket, server, cert, key string) {
	var ln net.Listener
	var err error
	if activatedFds := listenFds(); len(activatedFds) == 0 {
		ln, err = net.Listen("tcp", server)
		if err != nil {
			panic(err)
		}
	} else if len(activatedFds) == 1 {
		fmt.Println("socket activation!") //TODO: remove
		ln, err = net.FileListener(activatedFds[0])
		if err != nil {
			panic(err)
		}
		//TODO: does activatedFDs[0] need to be closed?
	} else {
		panic("Too many activated sockets! Check .socket file configuration.")
	}

	if key != "" {
		ln = wrapTLS(ln, cert, key)
	} else {
		fmt.Println("danger: no certificate or key specified - starting without TLS!")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "accept failure:", err)
			continue
		}

		unix, err := net.Dial("unix", socket)
		if err != nil {
			fmt.Fprintln(os.Stderr, "socket connection error:", err)
			conn.Close()
			continue
		}

		go func() {
			//TODO: consider ways the connection might close
			go io.Copy(unix, conn)
			io.Copy(conn, unix)
			conn.Close()
		}()
	}
}

func wrapTLS(listener net.Listener, certFile, keyFile string) net.Listener {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "certificate load failed:", err)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		/* ServerName: "test_server", */
		/* ClientAuth: tls.RequireAndVerifyClientCert, */
		ClientAuth: tls.RequireAnyClientCert,
	}
	cfg.BuildNameToCertificate()

	return tls.NewListener(listener, cfg)
}

// based on: https://gist.github.com/alberts/4640792
const (
	listenFdsStart = 3
)

func listenFds() []*os.File {
	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil || pid != os.Getpid() {
		return nil
	}
	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds == 0 {
		return nil
	}
	files := []*os.File(nil)
	for fd := listenFdsStart; fd < listenFdsStart+nfds; fd++ {
		syscall.CloseOnExec(fd)
		files = append(files, os.NewFile(uintptr(fd), ""))
	}
	return files
}
