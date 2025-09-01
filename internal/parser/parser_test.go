package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/token"
)

func stringPtr(s string) *string {
	return &s
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  token.Item
		expectErr bool
	}{
		{
			name:     "Simple String",
			input:    "+OK\r\n",
			expected: &token.SimpleString{Value: "OK"},
		},
		{
			name:     "Error",
			input:    "-Error message\r\n",
			expected: &token.Error{Value: "Error message"},
		},
		{
			name:     "Positive Integer",
			input:    ":1000\r\n",
			expected: &token.Integer{Value: 1000},
		},
		{
			name:     "Negative Integer",
			input:    ":-25\r\n",
			expected: &token.Integer{Value: -25},
		},
		{
			name:     "Bulk String",
			input:    "$5\r\nhello\r\n",
			expected: &token.BulkString{Value: stringPtr("hello")},
		},
		{
			name:     "Empty Bulk String",
			input:    "$0\r\n\r\n",
			expected: &token.BulkString{Value: stringPtr("")},
		},
		{
			name:     "Null Bulk String",
			input:    "$-1\r\n",
			expected: &token.BulkString{Value: nil},
		},
		{
			name:     "Bulk String with CRLF inside",
			input:    "$10\r\nwassup\r\nAB\r\n",
			expected: &token.BulkString{Value: stringPtr("wassup\r\nAB")},
		},
		{
			name:     "Empty Array",
			input:    "*0\r\n",
			expected: &token.Array{Items: []token.Item{}},
		},
		{
			name:     "Null Array",
			input:    "*-1\r\n",
			expected: &token.Array{Items: []token.Item{}},
		},
		{
			name:  "Array of Bulk Strings",
			input: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			expected: &token.Array{
				Items: []token.Item{
					&token.BulkString{Value: stringPtr("hello")},
					&token.BulkString{Value: stringPtr("world")},
				},
			},
		},
		{
			name:  "Array of Mixed Types",
			input: "*3\r\n:1\r\n+OK\r\n$5\r\nhello\r\n",
			expected: &token.Array{
				Items: []token.Item{
					&token.Integer{Value: 1},
					&token.SimpleString{Value: "OK"},
					&token.BulkString{Value: stringPtr("hello")},
				},
			},
		},
		{
			name:  "Nested Array",
			input: "*2\r\n*2\r\n:1\r\n:2\r\n*1\r\n+Nested\r\n",
			expected: &token.Array{
				Items: []token.Item{
					&token.Array{
						Items: []token.Item{
							&token.Integer{Value: 1},
							&token.Integer{Value: 2},
						},
					},
					&token.Array{
						Items: []token.Item{
							&token.SimpleString{Value: "Nested"},
						},
					},
				},
			},
		},

		// --- Error Cases ---
		{
			name:      "Invalid Type Prefix",
			input:     "!invalid\r\n",
			expectErr: true,
		},
		{
			name:      "Malformed Integer",
			input:     ":notanumber\r\n",
			expectErr: true,
		},
		{
			name:      "Malformed Bulk String Length",
			input:     "$foo\r\nbar\r\n",
			expectErr: true,
		},
		{
			name:      "Incomplete Bulk String Data",
			input:     "$10\r\nhello\r\n", // Requesting 10 bytes, but only 5 are provided + newline
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			p := NewParser(reader)

			item, err := p.Parse()

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected an error, but got none")
				}
				return // Test is done if an error was correctly expected.
			}

			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			if !reflect.DeepEqual(item, tc.expected) {
				t.Errorf("parsed item does not match expected item.\n got: %#v\nwant: %#v", item, tc.expected)
			}
		})
	}
}
