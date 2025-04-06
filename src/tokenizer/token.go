package tokenizer

type TokenType int

const (
	Number TokenType = iota
	Operator
	Function
	LeftBrace
	RightBrace
	Comma
	Constant
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
	"^": Operator,
}

var Functions = map[string]TokenType{
	"sin":   Function,
	"cos":   Function,
	"tg":    Function,
	"ctg":   Function,
	"log":   Function,
	"log2":  Function,
	"log10": Function,
	"sqrt":  Function,
	"abs":   Function,
}

var Constants = map[string]TokenType{
	"pi": Constant,
	"e":  Constant,
}
