package tokenizer

import (
	"unicode"
)

func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	runes := []rune(input)
	i := 0
	var prevToken Token

	for i < len(runes) {
		r := runes[i]

		switch {
		case unicode.IsSpace(r):
			i++

		case r == '-' && (i == 0 || isOperator(prevToken)):
			start := i
			i++
			if i >= len(runes) || !unicode.IsDigit(runes[i]) {
				return nil, ErrInvalidNumber(start)
			}
			for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{Number, string(runes[start:i]), start})
			prevToken = tokens[len(tokens)-1]

		case unicode.IsDigit(r) || r == '.':
			start := i
			dotCount := 0
			if r == '.' {
				dotCount++
			}

			for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				if runes[i] == '.' {
					dotCount++
					if dotCount > 1 {
						return nil, ErrInvalidNumber(i)
					}
				}
				i++
			}
			if dotCount == 1 && (start == i-1 || runes[start] == '.') {
				return nil, ErrInvalidNumber(start)
			}

			if i < len(runes) && (runes[i] == 'e' || runes[i] == 'E') {
				ePos := i
				i++

				if i < len(runes) && (runes[i] == '+' || runes[i] == '-') {
					i++
				}

				if i >= len(runes) || !unicode.IsDigit(runes[i]) {
					return nil, ErrInvalidNumber(ePos)
				}

				for i < len(runes) && unicode.IsDigit(runes[i]) {
					i++
				}

				if i < len(runes) && unicode.IsLetter(runes[i]) {
					return nil, ErrInvalidNumber(i)
				}
			}

			tokens = append(tokens, Token{Number, string(runes[start:i]), start})
			prevToken = tokens[len(tokens)-1]

		case unicode.IsLetter(r):
			start := i
			for i < len(runes) && unicode.IsLetter(runes[i]) {
				i++
			}
			if i > 0 && start > 0 && unicode.IsDigit(runes[start-1]) {
				return nil, ErrInvalidNumber(start)
			}
			tokens = append(tokens, Token{Function, string(runes[start:i]), start})
			prevToken = tokens[len(tokens)-1]

		case isSupportedOperator(r):
			if len(tokens) > 0 && isOperator(prevToken) {
				return nil, ErrInvalidRPNSyntax(i)
			}
			tokens = append(tokens, Token{Operator, string(r), i})
			prevToken = tokens[len(tokens)-1]
			i++

		case r == '(' || r == ')':
			tokType := LeftBrace
			if r == ')' {
				tokType = RightBrace
			}
			tokens = append(tokens, Token{tokType, string(r), i})
			prevToken = tokens[len(tokens)-1]
			i++

		default:
			return nil, ErrUnknownSymbol(i)
		}
	}

	if err := validateExpressionStructure(tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func validateExpressionStructure(tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}

	first := tokens[0].Type
	last := tokens[len(tokens)-1].Type

	if first == Operator || first == RightBrace {
		return ErrInvalidRPNSyntax(tokens[0].Pos)
	}

	if last == Operator || last == LeftBrace {
		return ErrInvalidRPNSyntax(tokens[len(tokens)-1].Pos)
	}

	parenCount := 0
	for i, token := range tokens {
		switch token.Type {
		case Number:
			if i > 0 && tokens[i-1].Type == Number {
				return ErrInvalidNumber(token.Pos)
			}
		case LeftBrace:
			parenCount++
			if i > 0 && tokens[i-1].Type == Number {
				return ErrNotEnoughOperands
			}
		case RightBrace:
			parenCount--
			if parenCount < 0 {
				return ErrMismatchedParentheses(token.Pos)
			}
			if i > 0 && tokens[i-1].Type == Operator {
				return ErrNotEnoughOperands
			}
		case Operator:
			if i == 0 || i == len(tokens)-1 {
				continue
			}
			prev := tokens[i-1].Type
			next := tokens[i+1].Type
			if prev == Operator || prev == LeftBrace {
				return ErrInvalidRPNSyntax(token.Pos)
			}
			if next == Operator || next == RightBrace {
				return ErrInvalidRPNSyntax(token.Pos)
			}
		}
	}

	if parenCount != 0 {
		return ErrMismatchedParentheses(0)
	}

	return nil
}

func isSupportedOperator(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '^':
		return true
	default:
		return false
	}
}

func isOperator(t Token) bool {
	if t.Type != Operator {
		return false
	}
	return isSupportedOperator(rune(t.Value[0]))
}
