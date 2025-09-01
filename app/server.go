package main

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/token"
)

func main() {
	command := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	fmt.Printf("--- Parsing command ---\n%s\n", command)
	reader := strings.NewReader(command)
	p := parser.NewParser(reader)

	item, err := p.Parse()
	if err != nil {
		fmt.Println("Error parsing command:", err)
		return
	}
	inspectItem(item)
}

func inspectItem(item token.Item) {
	switch v := item.(type) {
	case *token.Array:
		fmt.Printf("Parsed an Array with %d elements:\n", len(v.Items))
		for i, subItem := range v.Items {
			fmt.Printf("  [%d] -> ", i)
			inspectItem(subItem)
		}
	case *token.BulkString:
		if v.Value == nil {
			fmt.Println("Parsed a Null Bulk String.")
		} else {
			fmt.Printf("Parsed a Bulk String: \"%s\"\n", *v.Value)
		}
	case *token.SimpleString:
		fmt.Printf("Parsed a Simple String: \"%s\"\n", v.Value)
	case *token.Integer:
		fmt.Printf("Parsed an Integer: %d\n", v.Value)
	case *token.Error:
		fmt.Printf("Parsed an Error: \"%s\"\n", v.Value)
	default:
		fmt.Printf("Parsed an unknown item of type %T\n", v)
	}
}
