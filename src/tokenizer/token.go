package tokenizer

type TokenType int

const (
	Number TokenType = iota
	Operator
	Function
	LeftBrace
	RightBrace
	Comma
)

type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

var Operators = map[string]TokenType{
	"+": Operator,
	"-": Operator,
	"*": Operator,
	"/": Operator,
}

var Functions = map[string]TokenType{}
