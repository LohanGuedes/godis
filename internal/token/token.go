package token

import "strconv"

type RESPType byte

const (
	SIMPLE_STRING RESPType = '+'
	ERROR         RESPType = '-'
	INTEGER       RESPType = ':'
	BULK_STRING   RESPType = '$'
	ARRAY         RESPType = '*'
	NULL          RESPType = '-'
)

type Item interface {
	Type() RESPType
	Literal() string
}

type SimpleString struct {
	Value string
}

func (s *SimpleString) Type() RESPType  { return SIMPLE_STRING }
func (s *SimpleString) Literal() string { return s.Value }

type Integer struct {
	Value int64
}

func (i *Integer) Type() RESPType  { return INTEGER }
func (i *Integer) Literal() string { return strconv.FormatInt(i.Value, 10) }

type Error struct {
	Value string
}

func (e *Error) Type() RESPType  { return ERROR }
func (e *Error) Literal() string { return e.Value }

type BulkString struct {
	Value *string
}

func (b *BulkString) Type() RESPType { return BULK_STRING }
func (b *BulkString) Literal() string {
	if b.Value == nil {
		return ""
	}
	return *b.Value
}

type Array struct {
	Items []Item
}

func (a *Array) Type() RESPType  { return ARRAY }
func (a *Array) Literal() string { return "array" }
