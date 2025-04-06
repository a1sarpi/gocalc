package evaluation

import (
	"github.com/a1sarpi/gocalc/src/stack"
	"github.com/a1sarpi/gocalc/src/tokenizer"
)

type Associativity int

const (
	LeftAssociativity Associativity = iota
	RightAssociativity
)

type Stack stack.Stack[tokenizer.Token]

func ToRPN(tokens []tokenizer.Token) ([]tokenizer.Token, error) {
	result := make([]tokenizer.Token, 0, len(tokens))
	operationStack := stack.New[tokenizer.Token]()

	for _, tok := range tokens {
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
	}

	for !operationStack.IsEmpty() {
		if operationStack.Top().Type == tokenizer.LeftBrace {
			return nil, tokenizer.ErrMismatchedParentheses(operationStack.Top().Pos)
		}

		result = append(result, operationStack.Pop())
	}

	return result, nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func associativity(op string) Associativity {
	return LeftAssociativity
}
