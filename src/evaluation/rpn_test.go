package evaluation_test

import (
	"github.com/a1sarpi/gocalc/src/evaluation"
	"github.com/a1sarpi/gocalc/src/tokenizer"
	"testing"
)

func TestToRPN(t *testing.T) {
	tests := []struct {
		name  string
		infix []tokenizer.Token
		want  []tokenizer.Token
	}{
		{
			"Simple addition",
			[]tokenizer.Token{
				{tokenizer.Number, "3", 0},
				{tokenizer.Operator, "+", 1},
				{tokenizer.Number, "4", 2},
			},
			[]tokenizer.Token{
				{tokenizer.Number, "3", 0},
				{tokenizer.Number, "4", 2},
				{tokenizer.Operator, "+", 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluation.ToRPN(tt.infix)
			if err != nil {
				t.Fatalf("ToRPN() error = %v", err)
			}
			if !compareTokenSlices(got, tt.want) {
				t.Errorf("ToRPN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareTokenSlices(a, b []tokenizer.Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Type != b[i].Type {
			return false
		}
		if a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}
