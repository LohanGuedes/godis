package eventloop

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/token"
)

// Start begins listening for and handling TCP connections on port 6379.
func Start() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379:", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Server is listening on port 6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

// handleConnection manages a single client connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create the parser ONCE, outside the loop.
	// This ensures the same buffer is used for the entire connection.
	p := parser.NewParser(conn)

	// Loop to continuously read and process commands from the client.
	for {
		parsedItem, err := p.Parse()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected.")
				break
			}
			fmt.Printf("Error parsing command: %v\n", err)
			return
		}

		cmdArray, ok := parsedItem.(*token.Array)
		if !ok {
			conn.Write([]byte("-ERR Commands must be sent as arrays\r\n"))
			continue
		}
		if len(cmdArray.Items) == 0 {
			conn.Write([]byte("-ERR Empty command\r\n"))
			continue
		}

		commandNameItem, ok := cmdArray.Items[0].(*token.BulkString)
		if !ok || commandNameItem.Value == nil {
			conn.Write([]byte("-ERR Invalid command format\r\n"))
			continue
		}
		command := strings.ToUpper(*commandNameItem.Value)

		switch command {
		case "PING":
			handlePing(conn, cmdArray.Items)
		case "ECHO":
			handleEcho(conn, cmdArray.Items)
		default:
			errMsg := fmt.Sprintf("-ERR unknown command `%s`\r\n", *commandNameItem.Value)
			conn.Write([]byte(errMsg))
		}
	}
}

// handlePing responds to a PING command.
func handlePing(conn net.Conn, args []token.Item) {
	if len(args) > 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ping' command\r\n"))
		return
	}

	if len(args) == 2 {
		arg, ok := args[1].(*token.BulkString)
		if !ok || arg.Value == nil {
			conn.Write([]byte("-ERR PING argument must be a Bulk String\r\n"))
			return
		}
		response := fmt.Sprintf("$%d\r\n%s\r\n", len(*arg.Value), *arg.Value)
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte("+PONG\r\n"))
	}
}

// handleEcho responds to an ECHO command.
func handleEcho(conn net.Conn, args []token.Item) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'echo' command\r\n"))
		return
	}

	arg, ok := args[1].(*token.BulkString)
	if !ok || arg.Value == nil {
		conn.Write([]byte("-ERR ECHO argument must be a bulk string\r\n"))
		return
	}

	response := fmt.Sprintf("$%d\r\n%s\r\n", len(*arg.Value), *arg.Value)
	conn.Write([]byte(response))
}
