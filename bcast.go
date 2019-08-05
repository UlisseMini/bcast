package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func run(nConns int, l net.Listener, stream io.Reader) {
	conns := make([]io.Writer, nConns)
	fmt.Println("Waiting for connections...")
	for i := 1; i <= nConns; i++ {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		fmt.Printf("Connection from %s (need %d more)\n", conn.RemoteAddr(), nConns-i)
		conns[i-1] = conn
		defer func() {
			if err := conn.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
	}

	fmt.Printf("Streaming to all %d connected devices\n", nConns)
	multi := io.MultiWriter(conns...)

	_, err := io.Copy(multi, stream)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
	nConns := flag.Int("m", 1, "amount of connections to accept before streaming")
	addr := flag.String("l", ":4444", "address to listen on")
	proto := flag.String("p", "tcp", "protocol, can be [udp, tcp, udp6, tcp6] etc")

	flag.Parse()

	listener, err := net.Listen(*proto, *addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	run(*nConns, listener, os.Stdin)
}
