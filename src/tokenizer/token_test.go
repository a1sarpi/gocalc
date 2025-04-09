package tokenizer_test

import (
	"testing"

	"github.com/a1sarpi/gocalc/src/tokenizer"
)

func TestTokenizeNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokenizer.Token
	}{
		{
			name:  "single digit",
			input: "5",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "5", Pos: 0},
			},
		},
		{
			name:  "two digits",
			input: "42",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "42", Pos: 0},
			},
		},
		{
			name:  "three digits",
			input: "123",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "123", Pos: 0},
			},
		},
		{
			name:  "decimal number",
			input: "3.14",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "3.14", Pos: 0},
			},
		},
		{
			name:  "negative number",
			input: "-42",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "-42", Pos: 0},
			},
		},
		{
			name:  "scientific notation positive exponent",
			input: "1.25e+09",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "1.25e+09", Pos: 0},
			},
		},
		{
			name:  "scientific notation negative exponent",
			input: "1.25e-09",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "1.25e-09", Pos: 0},
			},
		},
		{
			name:  "scientific notation no sign",
			input: "1e5",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "1e5", Pos: 0},
			},
		},
		{
			name:  "scientific notation capital E",
			input: "1.25E+09",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "1.25E+09", Pos: 0},
			},
		},
		{
			name:  "unary minus at start",
			input: "-5",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "-5", Pos: 0},
			},
		},
		{
			name:  "unary minus after operator",
			input: "2 + -5",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "2", Pos: 0},
				{Type: tokenizer.Operator, Value: "+", Pos: 2},
				{Type: tokenizer.Number, Value: "-5", Pos: 4},
			},
		},
		{
			name:  "unary minus after parenthesis",
			input: "(-5)",
			want: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "(", Pos: 0},
				{Type: tokenizer.Number, Value: "-5", Pos: 1},
				{Type: tokenizer.RightBrace, Value: ")", Pos: 3},
			},
		},
		{
			name:  "unary minus in expression",
			input: "2 * (-5)",
			want: []tokenizer.Token{
				{Type: tokenizer.Number, Value: "2", Pos: 0},
				{Type: tokenizer.Operator, Value: "*", Pos: 2},
				{Type: tokenizer.LeftBrace, Value: "(", Pos: 4},
				{Type: tokenizer.Number, Value: "-5", Pos: 5},
				{Type: tokenizer.RightBrace, Value: ")", Pos: 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize(%q) failed: %v", tt.input, err)
			}
			if !compareTokens(got, tt.want) {
				t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTokenizeOperations(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokenizer.Token
	}{
		{
			name:  "addition",
			input: "1 + 2",
			want: []tokenizer.Token{
				{tokenizer.Number, "1", 0},
				{tokenizer.Operator, "+", 2},
				{tokenizer.Number, "2", 4},
			},
		},
		{
			name:  "subtraction",
			input: "5 - 3",
			want: []tokenizer.Token{
				{tokenizer.Number, "5", 0},
				{tokenizer.Operator, "-", 2},
				{tokenizer.Number, "3", 4},
			},
		},
		{
			name:  "multiplication",
			input: "4 * 2",
			want: []tokenizer.Token{
				{tokenizer.Number, "4", 0},
				{tokenizer.Operator, "*", 2},
				{tokenizer.Number, "2", 4},
			},
		},
		{
			name:  "division",
			input: "8 / 4",
			want: []tokenizer.Token{
				{tokenizer.Number, "8", 0},
				{tokenizer.Operator, "/", 2},
				{tokenizer.Number, "4", 4},
			},
		},
		{
			name:  "multiple operations",
			input: "1 + 2 * 3 - 4 / 5",
			want: []tokenizer.Token{
				{tokenizer.Number, "1", 0},
				{tokenizer.Operator, "+", 2},
				{tokenizer.Number, "2", 4},
				{tokenizer.Operator, "*", 6},
				{tokenizer.Number, "3", 8},
				{tokenizer.Operator, "-", 10},
				{tokenizer.Number, "4", 12},
				{tokenizer.Operator, "/", 14},
				{tokenizer.Number, "5", 16},
			},
		},
		{
			name:  "no spaces",
			input: "1+2",
			want: []tokenizer.Token{
				{tokenizer.Number, "1", 0},
				{tokenizer.Operator, "+", 1},
				{tokenizer.Number, "2", 2},
			},
		},
		{
			name:  "exponentiation",
			input: "2 ^ 3",
			want: []tokenizer.Token{
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "^", 2},
				{tokenizer.Number, "3", 4},
			},
		},
		{
			name:  "parentheses",
			input: "(2 + 3) * 4",
			want: []tokenizer.Token{
				{tokenizer.LeftBrace, "(", 0},
				{tokenizer.Number, "2", 1},
				{tokenizer.Operator, "+", 3},
				{tokenizer.Number, "3", 5},
				{tokenizer.RightBrace, ")", 6},
				{tokenizer.Operator, "*", 8},
				{tokenizer.Number, "4", 10},
			},
		},
		{
			name:  "nested parentheses",
			input: "2 * (3 + (4 - 1))",
			want: []tokenizer.Token{
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "*", 2},
				{tokenizer.LeftBrace, "(", 4},
				{tokenizer.Number, "3", 5},
				{tokenizer.Operator, "+", 7},
				{tokenizer.LeftBrace, "(", 9},
				{tokenizer.Number, "4", 10},
				{tokenizer.Operator, "-", 12},
				{tokenizer.Number, "1", 14},
				{tokenizer.RightBrace, ")", 15},
				{tokenizer.RightBrace, ")", 16},
			},
		},
		{
			name:  "exponentiation with parentheses",
			input: "2 ^ (3 + 1)",
			want: []tokenizer.Token{
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "^", 2},
				{tokenizer.LeftBrace, "(", 4},
				{tokenizer.Number, "3", 5},
				{tokenizer.Operator, "+", 7},
				{tokenizer.Number, "1", 9},
				{tokenizer.RightBrace, ")", 10},
			},
		},
		{
			name:  "pi constant",
			input: "pi",
			want: []tokenizer.Token{
				{tokenizer.Constant, "pi", 0},
			},
		},
		{
			name:  "e constant",
			input: "e",
			want: []tokenizer.Token{
				{tokenizer.Constant, "e", 0},
			},
		},
		{
			name:  "expression with constants",
			input: "pi * e",
			want: []tokenizer.Token{
				{tokenizer.Constant, "pi", 0},
				{tokenizer.Operator, "*", 3},
				{tokenizer.Constant, "e", 5},
			},
		},
		{
			name:  "function with constant",
			input: "sin(pi)",
			want: []tokenizer.Token{
				{tokenizer.Function, "sin", 0},
				{tokenizer.LeftBrace, "(", 3},
				{tokenizer.Constant, "pi", 4},
				{tokenizer.RightBrace, ")", 6},
			},
		},
		{
			name:  "unary minus before function",
			input: "-sin(1)",
			want: []tokenizer.Token{
				{tokenizer.Operator, "-", 0},
				{tokenizer.Function, "sin", 1},
				{tokenizer.LeftBrace, "(", 4},
				{tokenizer.Number, "1", 5},
				{tokenizer.RightBrace, ")", 6},
			},
		},
		{
			name:  "unary minus before parentheses",
			input: "-(1 + 2)",
			want: []tokenizer.Token{
				{tokenizer.Operator, "-", 0},
				{tokenizer.LeftBrace, "(", 1},
				{tokenizer.Number, "1", 2},
				{tokenizer.Operator, "+", 4},
				{tokenizer.Number, "2", 6},
				{tokenizer.RightBrace, ")", 7},
			},
		},
		{
			name:  "unary minus before function in expression",
			input: "2 * -sin(1)",
			want: []tokenizer.Token{
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "*", 2},
				{tokenizer.Operator, "-", 4},
				{tokenizer.Function, "sin", 5},
				{tokenizer.LeftBrace, "(", 8},
				{tokenizer.Number, "1", 9},
				{tokenizer.RightBrace, ")", 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize(%q) failed: %v", tt.input, err)
			}
			if !compareTokens(got, tt.want) {
				t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "incomplete operation", input: "2 /"},
		{name: "invalid number", input: "1.2.3"},
		{name: "double operator", input: "1 ++ 2"},
		{name: "missing operand", input: "1 + * 2"},
		{name: "invalid characters", input: "1 @ 2"},
		{name: "space in number", input: "1 1 + 1"},
		{name: "complex number", input: "1 + 4j"},
		{name: "incomplete scientific notation", input: "1.2e"},
		{name: "scientific notation no exponent", input: "1.2e+"},
		{name: "scientific notation decimal exponent", input: "1.2e1.2"},
		{name: "scientific notation double sign", input: "1.2e++2"},
		{name: "scientific notation with letter", input: "1e10f"},
		{name: "mismatched parentheses", input: "(2 + 3))"},
		{name: "unclosed parentheses", input: "(2 + 3"},
		{name: "operator after right parenthesis", input: "(2 + 3) +"},
		{name: "operator before left parenthesis", input: "2 + (3"},
		{name: "unknown constant", input: "unknown"},
		{name: "constant with number", input: "pi2"},
		{name: "constant with letter", input: "pix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tokenizer.Tokenize(tt.input)
			if err == nil {
				t.Errorf("Expected error for input: %q", tt.input)
			}
		})
	}
}

func compareTokens(a, b []tokenizer.Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Type != b[i].Type || a[i].Value != b[i].Value || a[i].Pos != b[i].Pos {
			return false
		}
	}
	return true
}
