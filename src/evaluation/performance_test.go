package evaluation_test

import (
	"strings"
	"testing"
	"time"

	"github.com/a1sarpi/gocalc/src/evaluation"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

const maxAllowedTime = 200 * time.Millisecond

func TestPerformanceAndLimits(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedError      bool
		skipTimeLimitCheck bool
	}{
		{
			name:          "1000 ones addition",
			input:         strings.Repeat("1 + ", 999) + "1",
			expectedError: false,
		},
		{
			name:          "Large number addition",
			input:         "1" + strings.Repeat("0", 300) + " + " + "1" + strings.Repeat("0", 300),
			expectedError: false,
		},
		{
			name:          "Deeply nested parentheses",
			input:         strings.Repeat("(", 100) + "1" + strings.Repeat(")", 100),
			expectedError: false,
		},
		{
			name:          "Long trigonometric chain",
			input:         strings.Repeat("sin(", 50) + "1" + strings.Repeat(")", 50),
			expectedError: false,
		},
		{
			name:          "Complex exponential",
			input:         "1.000000000000001 ^ 36893488147419103232",
			expectedError: true,
		},
		{
			name:          "Large factorial simulation",
			input:         strings.Repeat("2 * ", 100) + "1",
			expectedError: false,
		},
		{
			name:          "Many decimal places",
			input:         "3." + strings.Repeat("1", 1000) + " + 2." + strings.Repeat("9", 1000),
			expectedError: false,
		},
		{
			name:          "Invalid long expression",
			input:         strings.Repeat("1 + ", 1000) + "+",
			expectedError: true,
		},
		{
			name:          "Long expression with constants",
			input:         strings.Repeat("pi + ", 500) + "e",
			expectedError: false,
		},
		{
			name:          "Complex logarithmic chain",
			input:         strings.Repeat("ln(", 5) + "e" + strings.Repeat(")", 5),
			expectedError: true,
		},
		{
			name:          "Mixed operations with large numbers",
			input:         "1e100 * 2e100 + 3e100 / 4e100 - 5e100",
			expectedError: false,
		},
		{
			name:               "Very long but valid expression",
			input:              strings.Repeat("(1 + 2 * 3) / 4 + ", 200) + "5",
			expectedError:      false,
			skipTimeLimitCheck: true,
		},
		{
			name:          "Extreme precision calculation",
			input:         "(" + strings.Repeat("1/2 + ", 300) + "1/2)",
			expectedError: false,
		},
		{
			name:          "Maximum precision trigonometry",
			input:         "sin(" + strings.Repeat("pi/6 + ", 300) + "pi/6)",
			expectedError: false,
		},
		{
			name:          "Overflow through multiplication",
			input:         strings.Repeat("1e10 * ", 50) + "1e10",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			tokens, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				if !tt.expectedError {
					t.Errorf("Unexpected tokenization error: %v", err)
				}
				return
			}

			rpn, err := evaluation.ToRPN(tokens)
			if err != nil {
				if !tt.expectedError {
					t.Errorf("Unexpected RPN conversion error: %v", err)
				}
				return
			}

			_, err = evaluation.Calculate(rpn, false)
			duration := time.Since(start)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.skipTimeLimitCheck && duration > maxAllowedTime {
				t.Errorf("Execution took too long: %v (limit: %v)", duration, maxAllowedTime)
			}

			t.Logf("Expression length: %d, Execution time: %v", len(tt.input), duration)
		})
	}
}

func TestExtremeInputs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Empty expression",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Single space",
			input:   " ",
			wantErr: true,
		},
		{
			name:    "Only operators",
			input:   "+-*/",
			wantErr: true,
		},
		{
			name:    "Invalid characters",
			input:   "1 + 2 @ 3",
			wantErr: true,
		},
		{
			name:    "Unicode characters",
			input:   "1 + 2 × 3 ÷ 4",
			wantErr: true,
		},
		{
			name:    "Maximum number length",
			input:   "1" + strings.Repeat("0", 100),
			wantErr: false,
		},
		{
			name:    "Extremely small number",
			input:   "1e-1000",
			wantErr: false,
		},
		{
			name:    "Extremely large number",
			input:   "1e1000",
			wantErr: true,
		},
		{
			name:    "Many decimal points",
			input:   "1" + strings.Repeat(".", 10) + "5",
			wantErr: true,
		},
		{
			name:    "Many parentheses",
			input:   strings.Repeat("(", 1000) + "1" + strings.Repeat(")", 999),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			tokens, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Unexpected tokenization error: %v", err)
				}
				return
			}

			_, err = evaluation.ToRPN(tokens)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Unexpected RPN conversion error: %v", err)
				}
				return
			}

			_, err = evaluation.Calculate(tokens, false)
			duration := time.Since(start)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if duration > maxAllowedTime {
				t.Errorf("Execution took too long: %v (limit: %v)", duration, maxAllowedTime)
			}

			t.Logf("Expression length: %d, Execution time: %v", len(tt.input), duration)
		})
	}
}
