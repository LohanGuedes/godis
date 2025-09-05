package godis

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/kv"
	"github.com/codecrafters-io/redis-starter-go/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/token"
)

type Server struct {
	kv     kv.Store
	config Config
	// TODO: Maybe add a collection of current connected clients
	// TODO: Add a logger
}

type Config struct {
	Host            string
	Port            int
	CleanerInterval time.Duration
}

func NewServer(c Config) *Server {
	return &Server{
		kv:     *kv.NewStore(c.CleanerInterval),
		config: c,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.config.Host+":"+strconv.Itoa(s.config.Port))
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Printf("Server is listening on port %d\n", s.config.Port)

	go s.kv.ExpiryHandler()

	for {
		// NOTE: Maybe a connection shuuld be more than only that. holding more info?
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	p := parser.NewParser(conn)

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
			s.handlePing(conn, cmdArray.Items)
		case "ECHO":
			s.handleEcho(conn, cmdArray.Items)
		case "SET":
			s.handleSet(conn, cmdArray.Items)
		case "GET":
			s.handleGet(conn, cmdArray.Items)
		default:
			errMsg := fmt.Sprintf("-ERR unknown command `%s`\r\n", *commandNameItem.Value)
			conn.Write([]byte(errMsg))
		}
	}
}

func (Server) handlePing(conn net.Conn, args []token.Item) {
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

func (Server) handleEcho(conn net.Conn, args []token.Item) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
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

func (s *Server) handleSet(conn net.Conn, args []token.Item) (err error) {
	var duration time.Duration
	if len(args) > 5 || len(args) < 3 {
		conn.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
	}

	key, ok := args[1].(*token.BulkString)
	if !ok || key.Literal() == "" {
		conn.Write([]byte("-ERR SET arguments must be in a bulk string\r\n"))
		return
	}
	val, ok := args[2].(*token.BulkString)
	if !ok || key.Literal() == "" {
		conn.Write([]byte("-ERR SET arguments must be in a bulk string\r\n"))
		return
	}

	if len(args) == 5 && strings.ToUpper(args[3].Literal()) == "PX" {
		duration, err = time.ParseDuration(args[4].Literal() + "ms")
		if err != nil {
			return err
		}
	}

	s.kv.Set(key.Literal(), val.Literal(), duration)
	conn.Write([]byte("+OK\r\n"))
	return err
}

func (s *Server) handleGet(conn net.Conn, args []token.Item) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
	}

	arg, ok := args[1].(*token.BulkString)
	if !ok || arg.Value == nil {
		conn.Write([]byte("-ERR GET arguments must be in a bulk string\r\n"))
		return
	}

	// Set in the KV databsase
	value, ok := s.kv.Get(args[1].Literal())
	if !ok {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	conn.Write(fmt.Appendf([]byte{}, "$%d\r\n%s\r\n", len(value), value))
}
