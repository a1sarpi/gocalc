package tokenizer

import "fmt"

type Error struct {
	Message string
	Pos     int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error while tokenizing (position: %d): %s", e.Pos, e.Message)
}

func NewError(msg string, pos int) *Error {
	return &Error{msg, pos}
}

var (
	ErrUnknownSymbol = func(pos int) *Error {
		return NewError("Unknown symbol", pos)
	}
	ErrInvalidNumber = func(pos int) *Error {
		return NewError("Invalid number", pos)
	}
)

var (
	ErrMismatchedParentheses = func(pos int) *Error {
		return NewError("Mismatched parentheses", pos)
	}
	ErrInvalidRPNSyntax = func(pos int) *Error {
		return NewError("Invalid RPN syntax", pos)
	}
)

var (
	ErrInvalidExpression = fmt.Errorf("Invalid expression")
	ErrDivisionByZero    = fmt.Errorf("division by zero")
	ErrNotEnoughOperands = fmt.Errorf("not enough operands")
	ErrUnknownOperator   = func(op string) error {
		return fmt.Errorf("unknown operator: %s", op)
	}
	ErrUnknownFunction = func(fn string) error {
		return fmt.Errorf("error in function %s", fn)
	}
)
