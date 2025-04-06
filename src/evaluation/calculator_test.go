package evaluation_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/a1sarpi/gocalc/src/evaluation"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		rpn       []tokenizer.Token
		want      float64
		wantError bool
	}{
		{
			name: "simple addition",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "3", 0},
				{tokenizer.Number, "4", 0},
				{tokenizer.Operator, "+", 0},
			},
			want:      7,
			wantError: false,
		},
		{
			name: "simple subtraction",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "5", 0},
				{tokenizer.Number, "3", 0},
				{tokenizer.Operator, "-", 0},
			},
			want:      2,
			wantError: false,
		},
		{
			name: "simple multiplication",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "4", 0},
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "*", 0},
			},
			want:      8,
			wantError: false,
		},
		{
			name: "simple division",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "8", 0},
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "/", 0},
			},
			want:      4,
			wantError: false,
		},
		{
			name: "division by zero",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "5", 0},
				{tokenizer.Number, "0", 0},
				{tokenizer.Operator, "/", 0},
			},
			wantError: true,
		},
		{
			name: "two operations - addition and multiplication",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "2", 0},
				{tokenizer.Number, "3", 0},
				{tokenizer.Number, "4", 0},
				{tokenizer.Operator, "*", 0},
				{tokenizer.Operator, "+", 0},
			},
			want:      14,
			wantError: false,
		},
		{
			name: "two operations - multiplication and division",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "6", 0},
				{tokenizer.Number, "2", 0},
				{tokenizer.Number, "3", 0},
				{tokenizer.Operator, "*", 0},
				{tokenizer.Operator, "/", 0},
			},
			want:      1,
			wantError: false,
		},
		{
			name: "two operations - subtraction and division",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "10", 0},
				{tokenizer.Number, "4", 0},
				{tokenizer.Operator, "-", 0},
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "/", 0},
			},
			want:      3,
			wantError: false,
		},
		{
			name: "floating point precision",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "1", 0},
				{tokenizer.Number, "0.0000000000000000000001", 0},
				{tokenizer.Operator, "+", 0},
			},
			want:      1,
			wantError: false,
		},
		{
			name: "arithmetic overflow",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "1e308", 0},
				{tokenizer.Number, "1e308", 0},
				{tokenizer.Operator, "*", 0},
			},
			wantError: true,
		},
		{
			name: "not enough operands",
			rpn: []tokenizer.Token{
				{tokenizer.Number, "2", 0},
				{tokenizer.Operator, "+", 0},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluation.Calculate(tt.rpn, false)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Calculate() unexpected error = %v", err)
			}

			if !almostEqual(got, tt.want) {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllOperators(t *testing.T) {
	ops := []struct {
		op   string
		a    float64
		b    float64
		want float64
	}{
		{"+", 5, 3, 8},
		{"-", 5, 3, 2},
		{"*", 5, 3, 15},
		{"/", 6, 3, 2},
	}

	for _, op := range ops {
		t.Run(op.op, func(t *testing.T) {
			rpn := []tokenizer.Token{
				{tokenizer.Number, floatToString(op.a), 0},
				{tokenizer.Number, floatToString(op.b), 0},
				{tokenizer.Operator, op.op, 0},
			}

			got, err := evaluation.Calculate(rpn, false)
			if err != nil {
				t.Errorf("Operator %q failed: %v", op.op, err)
				return
			}

			if !almostEqual(got, op.want) {
				t.Errorf("Operator %q = %v, want %v", op.op, got, op.want)
			}
		})
	}
}

func TestAllFunctions(t *testing.T) {
	tests := []struct {
		name    string
		fn      string
		arg     float64
		want    float64
		radians bool
	}{
		{"sin(30°)", "sin", 30, 0.5, false},
		{"cos(60°)", "cos", 60, 0.5, false},
		{"sin(π/2)", "sin", math.Pi / 2, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpn := []tokenizer.Token{
				{tokenizer.Number, floatToString(tt.arg), 0},
				{tokenizer.Function, tt.fn, 0},
			}

			got, err := evaluation.Calculate(rpn, tt.radians)
			if err != nil {
				t.Fatalf("Function %q failed: %v", tt.fn, err)
			}

			if !almostEqual(got, tt.want) {
				t.Errorf("Function %q = %v, want %v", tt.fn, got, tt.want)
			}
		})
	}
}

func TestCalculator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "simple addition",
			input:    "2 + 3",
			expected: 5,
		},
		{
			name:     "simple subtraction",
			input:    "5 - 3",
			expected: 2,
		},
		{
			name:     "simple multiplication",
			input:    "4 * 2",
			expected: 8,
		},
		{
			name:     "simple division",
			input:    "8 / 4",
			expected: 2,
		},
		{
			name:     "exponentiation",
			input:    "2 ^ 3",
			expected: 8,
		},
		{
			name:     "parentheses",
			input:    "(2 + 3) * 4",
			expected: 20,
		},
		{
			name:     "nested parentheses",
			input:    "2 * (3 + (4 - 1))",
			expected: 12,
		},
		{
			name:     "exponentiation with parentheses",
			input:    "2 ^ (3 + 1)",
			expected: 16,
		},
		{
			name:     "complex expression",
			input:    "2 + 3 * 4 ^ 2",
			expected: 50,
		},
		{
			name:     "scientific notation",
			input:    "1.23e5 + 4.56e-2",
			expected: 123000.0456,
		},
		{
			name:     "mixed operations",
			input:    "(2 + 3) * 4 ^ 2 - 10 / 2",
			expected: 75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize failed: %v", err)
			}

			rpn, err := evaluation.ToRPN(tokens)
			if err != nil {
				t.Fatalf("ToRPN failed: %v", err)
			}

			result, err := evaluation.Calculate(rpn, false)
			if err != nil {
				t.Fatalf("Calculate failed: %v", err)
			}

			if !almostEqual(result, tt.expected) {
				t.Errorf("Calculate(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	const epsilon = 1e-10
	return (a-b) < epsilon && (b-a) < epsilon
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
