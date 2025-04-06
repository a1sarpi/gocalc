package evaluation

import (
	"fmt"
	"math"
	"strconv"

	"github.com/a1sarpi/gocalc/src/stack"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

var ErrDivisionByZero = fmt.Errorf("division by zero")
var ErrInvalidRPNSyntax = fmt.Errorf("invalid RPN syntax")
var ErrArithmeticOverflow = fmt.Errorf("arithmetic overflow")

func ToRPN(tokens []tokenizer.Token) ([]tokenizer.Token, error) {
	output := make([]tokenizer.Token, 0, len(tokens))
	s := stack.New[tokenizer.Token]()

	for _, token := range tokens {
		switch token.Type {
		case tokenizer.Number, tokenizer.Constant:
			output = append(output, token)

		case tokenizer.Function:
			s.Push(token)

		case tokenizer.LeftBrace:
			s.Push(token)

		case tokenizer.RightBrace:
			for !s.IsEmpty() {
				top := s.Pop()
				if top.Type == tokenizer.LeftBrace {
					if !s.IsEmpty() {
						next := s.Top()
						if next.Type == tokenizer.Function {
							next = s.Pop()
							output = append(output, next)
						}
					}
					break
				}
				output = append(output, top)
			}

		case tokenizer.Operator:
			for !s.IsEmpty() {
				next := s.Top()
				if next.Type != tokenizer.Operator || !hasHigherPrecedence(next.Value, token.Value) {
					break
				}
				next = s.Pop()
				output = append(output, next)
			}
			s.Push(token)
		}
	}

	for !s.IsEmpty() {
		top := s.Pop()
		if top.Type == tokenizer.LeftBrace {
			return nil, tokenizer.ErrMismatchedParentheses(top.Pos)
		}
		output = append(output, top)
	}

	return output, nil
}

func Calculate(tokens []tokenizer.Token, useRadians bool) (float64, error) {
	s := stack.New[float64]()

	for _, token := range tokens {
		switch token.Type {
		case tokenizer.Number:
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, err
			}
			s.Push(val)

		case tokenizer.Constant:
			switch token.Value {
			case "pi":
				s.Push(math.Pi)
			case "e":
				s.Push(math.E)
			}

		case tokenizer.Function:
			if s.IsEmpty() {
				return 0, ErrInvalidRPNSyntax
			}

			x := s.Pop()

			switch token.Value {
			case "sin":
				if !useRadians {
					x = x * math.Pi / 180
				}
				s.Push(math.Sin(x))
			case "cos":
				if !useRadians {
					x = x * math.Pi / 180
				}
				s.Push(math.Cos(x))
			case "tan":
				if !useRadians {
					x = x * math.Pi / 180
				}
				s.Push(math.Tan(x))
			case "asin":
				result := math.Asin(x)
				if !useRadians {
					result = result * 180 / math.Pi
				}
				s.Push(result)
			case "acos":
				result := math.Acos(x)
				if !useRadians {
					result = result * 180 / math.Pi
				}
				s.Push(result)
			case "atan":
				result := math.Atan(x)
				if !useRadians {
					result = result * 180 / math.Pi
				}
				s.Push(result)
			case "log2":
				s.Push(math.Log2(x))
			case "log10":
				s.Push(math.Log10(x))
			case "sqrt":
				s.Push(math.Sqrt(x))
			case "abs":
				s.Push(math.Abs(x))
			}

		case tokenizer.Operator:
			if s.Len() < 2 {
				return 0, ErrInvalidRPNSyntax
			}

			b := s.Pop()
			a := s.Pop()

			switch token.Value {
			case "+":
				s.Push(a + b)
			case "-":
				s.Push(a - b)
			case "*":
				s.Push(a * b)
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				s.Push(a / b)
			case "^":
				s.Push(math.Pow(a, b))
			}
		}
	}

	if s.Len() != 1 {
		return 0, ErrInvalidRPNSyntax
	}

	result := s.Pop()
	return result, nil
}

func hasHigherPrecedence(op1, op2 string) bool {
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
	}
	return precedence[op1] >= precedence[op2]
}
