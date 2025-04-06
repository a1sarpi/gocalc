package evaluation

import (
	"fmt"
	"math"
	"strconv"

	"github.com/a1sarpi/gocalc/src/stack"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

type Associativity int

const (
	LeftAssociativity Associativity = iota
	RightAssociativity
)

var Debug bool

func ToRPN(tokens []tokenizer.Token) ([]tokenizer.Token, error) {
	result := make([]tokenizer.Token, 0, len(tokens))
	operationStack := stack.New[tokenizer.Token]()

	for _, tok := range tokens {
		if Debug {
			fmt.Printf("Processing token: %v\n", tok)
		}
		switch tok.Type {
		case tokenizer.Number:
			result = append(result, tok)

		case tokenizer.Operator:
			for !operationStack.IsEmpty() {
				top := operationStack.Top()
				if top.Type != tokenizer.Operator {
					break
				}
				if precedence(top.Value) < precedence(tok.Value) {
					break
				}
				result = append(result, operationStack.Pop())
			}
			operationStack.Push(tok)

		case tokenizer.Function:
			for !operationStack.IsEmpty() {
				top := operationStack.Top()
				if top.Type == tokenizer.Operator || top.Type == tokenizer.LeftBrace {
					break
				}
				if precedence(top.Value) > precedence(tok.Value) || (precedence(top.Value) == precedence(tok.Value) &&
					associativity(tok.Value) == LeftAssociativity) {
					result = append(result, operationStack.Pop())
					continue
				}
				break
			}
			operationStack.Push(tok)
		case tokenizer.LeftBrace:
			operationStack.Push(tok)
		case tokenizer.RightBrace:
			for !operationStack.IsEmpty() && operationStack.Top().Type != tokenizer.LeftBrace {
				result = append(result, operationStack.Pop())
			}
			if operationStack.IsEmpty() {
				return result, tokenizer.ErrMismatchedParentheses(tok.Pos)
			}
			operationStack.Pop()
			if !operationStack.IsEmpty() && operationStack.Top().Type == tokenizer.Function {
				result = append(result, operationStack.Pop())
			}
		default:
			return nil, tokenizer.ErrUnknownSymbol(tok.Pos)
		}
		if Debug {
			fmt.Printf("Current RPN: %v\n", result)
			fmt.Printf("Operation stack: %v\n", operationStack)
		}
	}

	for !operationStack.IsEmpty() {
		if operationStack.Top().Type == tokenizer.LeftBrace {
			return nil, tokenizer.ErrMismatchedParentheses(operationStack.Top().Pos)
		}

		result = append(result, operationStack.Pop())
	}

	if Debug {
		fmt.Printf("Final RPN: %v\n", result)
	}
	return result, nil
}

func Calculate(rpn []tokenizer.Token, useRadians bool) (float64, error) {
	numStack := stack.New[float64]()

	for _, token := range rpn {
		if Debug {
			fmt.Printf("Processing token: %v\n", token)
		}
		switch token.Type {
		case tokenizer.Number:
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				if Debug {
					fmt.Printf("Error parsing number %q: %v\n", token.Value, err)
				}
				return 0, tokenizer.ErrInvalidNumber(token.Pos)
			}
			numStack.Push(num)

		case tokenizer.Operator:
			if numStack.Len() < 2 {
				if Debug {
					fmt.Printf("Not enough operands for operator %q\n", token.Value)
				}
				return 0, tokenizer.ErrNotEnoughOperands
			}
			b := numStack.Pop()
			a := numStack.Pop()
			if Debug {
				fmt.Printf("Applying operator %q to %g and %g\n", token.Value, a, b)
			}
			res, err := applyOperator(token.Value, a, b)
			if err != nil {
				if Debug {
					fmt.Printf("Error applying operator: %v\n", err)
				}
				return 0, err
			}
			numStack.Push(res)

		case tokenizer.Function:
			if numStack.Len() < 1 {
				if Debug {
					fmt.Printf("Not enough operands for function %q\n", token.Value)
				}
				return 0, tokenizer.ErrNotEnoughOperands
			}
			arg := numStack.Pop()
			if Debug {
				fmt.Printf("Applying function %q to %g\n", token.Value, arg)
			}
			res, err := applyFunction(token.Value, arg, useRadians)
			if err != nil {
				if Debug {
					fmt.Printf("Error applying function: %v\n", err)
				}
				return 0, err
			}
			numStack.Push(res)

		default:
			if Debug {
				fmt.Printf("Unknown token type: %v\n", token.Type)
			}
			return 0, tokenizer.ErrUnknownSymbol(token.Pos)
		}
		if Debug {
			fmt.Printf("Number stack: %v\n", numStack)
		}
	}

	if numStack.Len() != 1 {
		if Debug {
			fmt.Printf("Invalid expression: stack has %d values\n", numStack.Len())
		}
		return 0, tokenizer.ErrInvalidExpression
	}

	result := numStack.Pop()
	if Debug {
		fmt.Printf("Final result: %g\n", result)
	}
	return result, nil
}

func applyOperator(op string, a, b float64) (float64, error) {
	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return 0, tokenizer.ErrDivisionByZero
		}
		result = a / b
	case "^":
		result = math.Pow(a, b)
	default:
		return 0, tokenizer.ErrUnknownOperator(op)
	}

	if math.IsInf(result, 0) {
		return 0, tokenizer.ErrInvalidExpression
	}

	return result, nil
}

func applyFunction(f string, a float64, useRadians bool) (float64, error) {
	if !useRadians && isTrig(f) {
		a = degreesToRadians(a)
	}

	switch f {
	case "sin":
		return math.Sin(a), nil

	case "cos":
		return math.Cos(a), nil

	case "tg":
		return math.Tan(a), nil

	case "ctg":
		return 1 / math.Tan(a), nil

	case "log":
		if a <= 0 {
			return 0, tokenizer.ErrInvalidExpression
		}
		return math.Log(a), nil

	case "sqrt":
		if a < 0 {
			return 0, tokenizer.ErrInvalidExpression
		}
		return math.Sqrt(a), nil

	default:
		return 0, tokenizer.ErrUnknownFunction(f)
	}
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	default:
		return 0
	}
}

func associativity(op string) Associativity {
	if op == "^" {
		return RightAssociativity
	}
	return LeftAssociativity
}

func isTrig(f string) bool {
	return f == "sin" || f == "cos" || f == "tg" || f == "ctg"
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}
