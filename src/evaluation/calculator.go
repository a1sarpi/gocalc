package evaluation

import (
	"math"
	"strconv"

	"github.com/a1sarpi/gocalc/src/stack"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

func Calculate(rpn []tokenizer.Token, useRadians bool) (float64, error) {
	numStack := stack.New[float64]()

	for _, token := range rpn {
		switch token.Type {
		case tokenizer.Number:
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, tokenizer.ErrInvalidNumber(token.Pos)
			}
			numStack.Push(num)

		case tokenizer.Operator:
			if numStack.Len() < 2 {
				return 0, tokenizer.ErrNotEnoughOperands
			}
			b := numStack.Pop()
			a := numStack.Pop()
			res, err := applyOperator(token.Value, a, b)
			if err != nil {
				return 0, err
			}
			numStack.Push(res)

		case tokenizer.Function:
			if numStack.Len() < 1 {
				return 0, tokenizer.ErrNotEnoughOperands
			}
			arg := numStack.Pop()
			res, err := applyFunction(token.Value, arg, useRadians)
			if err != nil {
				return 0, err
			}
			numStack.Push(res)

		default:
			return 0, tokenizer.ErrUnknownSymbol(token.Pos)
		}
	}

	if numStack.Len() != 1 {
		return 0, tokenizer.ErrInvalidExpression
	}

	return numStack.Pop(), nil
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

func isTrig(f string) bool {
	return f == "sin" || f == "cos" || f == "tg" || f == "ctg"
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}
