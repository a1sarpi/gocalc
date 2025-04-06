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
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i])) {
				i++
			}
			if i > 0 && start > 0 && unicode.IsDigit(runes[start-1]) {
				return nil, ErrInvalidNumber(start)
			}
			name := string(runes[start:i])
			if _, ok := Constants[name]; ok {
				tokens = append(tokens, Token{Constant, name, start})
			} else if _, ok := Functions[name]; ok {
				tokens = append(tokens, Token{Function, name, start})
			} else {
				return nil, ErrUnknownSymbol(start)
			}
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
		return ErrInvalidRPNSyntax(0)
	}

	if tokens[0].Type == Operator && tokens[0].Value != "-" {
		return ErrInvalidRPNSyntax(tokens[0].Pos)
	}

	if tokens[len(tokens)-1].Type == Operator {
		return ErrInvalidRPNSyntax(tokens[len(tokens)-1].Pos)
	}

	parenCount := 0
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case Operator:
			if i == len(tokens)-1 {
				return ErrInvalidRPNSyntax(token.Pos)
			}
			next := tokens[i+1]
			if next.Type != Number && next.Type != Constant && next.Type != LeftBrace && next.Type != Function && (next.Type != Operator || next.Value != "-") {
				return ErrInvalidRPNSyntax(token.Pos)
			}

		case Function:
			if i == len(tokens)-1 || tokens[i+1].Type != LeftBrace {
				return ErrInvalidRPNSyntax(token.Pos)
			}

		case LeftBrace:
			parenCount++
			if i == len(tokens)-1 {
				return ErrMismatchedParentheses(token.Pos)
			}
			next := tokens[i+1]
			if next.Type != Number && next.Type != Constant && next.Type != LeftBrace && next.Type != Function && (next.Type != Operator || next.Value != "-") {
				return ErrInvalidRPNSyntax(token.Pos)
			}

		case RightBrace:
			parenCount--
			if parenCount < 0 {
				return ErrMismatchedParentheses(token.Pos)
			}
			if i < len(tokens)-1 {
				next := tokens[i+1]
				if next.Type != Operator && next.Type != RightBrace {
					return ErrInvalidRPNSyntax(token.Pos)
				}
			}

		case Number, Constant:
			if i < len(tokens)-1 {
				next := tokens[i+1]
				if next.Type != Operator && next.Type != RightBrace {
					return ErrInvalidRPNSyntax(token.Pos)
				}
			}
		}

		if i == len(tokens)-1 && parenCount > 0 {
			return ErrMismatchedParentheses(token.Pos)
		}
	}

	if parenCount != 0 {
		return ErrMismatchedParentheses(tokens[len(tokens)-1].Pos)
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
