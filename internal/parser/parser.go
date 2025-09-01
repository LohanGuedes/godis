package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/token"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

func (p *Parser) readLine() ([]byte, error) {
	line, err := p.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return line[:len(line)-2], nil
}

func (p *Parser) Parse() (token.Item, error) {
	prefix, err := p.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch token.RESPType(prefix) {
	case token.SIMPLE_STRING:
		return p.parseSimpleString()
	case token.ERROR:
		return p.parseError()
	case token.INTEGER:
		return p.parseInteger()
	case token.BULK_STRING:
		return p.parseBulkString()
	case token.ARRAY:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("invalid or unsupported RESP type prefix: %q", prefix)
	}
}

func (p *Parser) parseSimpleString() (token.Item, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	return &token.SimpleString{Value: string(line)}, nil
}

func (p *Parser) parseError() (*token.Error, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	return &token.Error{Value: string(line)}, nil
}

func (p *Parser) parseInteger() (*token.Integer, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	val, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse integer: %w", err)
	}
	return &token.Integer{Value: val}, nil
}

func (p *Parser) parseBulkString() (*token.BulkString, error) {
	lengthLine, err := p.readLine()
	if err != nil {
		return nil, err
	}
	length, err := strconv.ParseInt(string(lengthLine), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse bulk string length: %w", err)
	}

	if length == -1 {
		return &token.BulkString{Value: nil}, nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(p.reader, data)
	if err != nil {
		return nil, err
	}

	if _, err = p.readLine(); err != nil {
		return nil, err
	}

	valStr := string(data)
	return &token.BulkString{Value: &valStr}, nil
}

// recursive
func (p *Parser) parseArray() (*token.Array, error) {
	countLine, err := p.readLine()
	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(string(countLine))
	if err != nil {
		return nil, fmt.Errorf("could not parse array items length: %w", err)
	}

	if count <= 0 {
		return &token.Array{Items: []token.Item{}}, nil
	}

	items := []token.Item{}
	for range count {
		item, err := p.Parse()
		if err != nil {
			return nil, err
		}
		items = append(items, item)

	}

	return &token.Array{Items: items}, nil
}
