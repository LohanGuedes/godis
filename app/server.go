package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(context.Background(), conn)
	}
}

func handleConn(ctx context.Context, c net.Conn) {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "PING" {
			c.Write([]byte("+PONG\r\n"))
		}
	}
}
