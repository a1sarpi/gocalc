package evaluation

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/a1sarpi/gocalc/src/stack"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

var (
	ErrDivisionByZero     = fmt.Errorf("division by zero")
	ErrInvalidRPNSyntax   = fmt.Errorf("invalid RPN syntax")
	ErrArithmeticOverflow = fmt.Errorf("arithmetic overflow")
	ErrTimeout            = fmt.Errorf("calculation timeout")
)

const (
	DefaultCalculationTime = 5 * time.Second
	TestCalculationTime    = 100 * time.Millisecond
	MaxFloat64             = 1.7976931348623157e+308
	MinFloat64             = -1.7976931348623157e+308
)

func checkOverflow(x float64) error {
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return ErrArithmeticOverflow
	}
	if x > MaxFloat64 || x < MinFloat64 {
		return ErrArithmeticOverflow
	}
	return nil
}

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
	return CalculateWithTimeout(tokens, useRadians, DefaultCalculationTime)
}

func CalculateWithTimeout(tokens []tokenizer.Token, useRadians bool, timeout time.Duration) (float64, error) {
	startTime := time.Now()
	s := stack.New[float64]()

	for _, token := range tokens {
		if time.Since(startTime) > timeout {
			return 0, ErrTimeout
		}

		switch token.Type {
		case tokenizer.Number:
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, err
			}
			if err := checkOverflow(val); err != nil {
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

			var result float64
			switch token.Value {
			case "sin":
				if !useRadians {
					x = x * math.Pi / 180
				}
				result = math.Sin(x)
			case "cos":
				if !useRadians {
					x = x * math.Pi / 180
				}
				result = math.Cos(x)
			case "tan":
				if !useRadians {
					x = x * math.Pi / 180
				}
				result = math.Tan(x)
			case "asin":
				result = math.Asin(x)
				if !useRadians {
					result = result * 180 / math.Pi
				}
			case "acos":
				result = math.Acos(x)
				if !useRadians {
					result = result * 180 / math.Pi
				}
			case "atan":
				result = math.Atan(x)
				if !useRadians {
					result = result * 180 / math.Pi
				}
			case "log":
				result = math.Log(x)
			case "log2":
				result = math.Log2(x)
			case "log10":
				result = math.Log10(x)
			case "sqrt":
				result = math.Sqrt(x)
			case "abs":
				result = math.Abs(x)
			case "exp":
				result = math.Exp(x)
			}
			if err := checkOverflow(result); err != nil {
				return 0, err
			}
			s.Push(result)

		case tokenizer.Operator:
			if s.Len() < 2 {
				return 0, ErrInvalidRPNSyntax
			}

			b := s.Pop()
			a := s.Pop()

			var result float64
			switch token.Value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				result = a / b
			case "^":
				result = math.Pow(a, b)
			}
			if err := checkOverflow(result); err != nil {
				return 0, err
			}
			s.Push(result)
		}
	}

	if s.Len() != 1 {
		return 0, ErrInvalidRPNSyntax
	}

	result := s.Pop()
	if err := checkOverflow(result); err != nil {
		return 0, err
	}
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
