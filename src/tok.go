package src

const (
	LEFT_CURLY = 0x1
	RIGHT_CURLY = 0x2
	LEFT_BRACKET = 0x3
	RIGHT_BRACKET = 0x4
	LEFT_PAREN = 0x5
	RIGHT_PAREN = 0x6
	COMMA = 0x7
	PERIOD = 0x8
	SEMICOLON = 0x9
	COLON = 0xA
	QUESTION = 0xB
	STAR = 0xC
	BANG = 0xD
	PERCENT = 0xE
	BIN_AND = 0xF
	BIN_OR = 0x10
	WIGGLE = 0x11
	PLUS = 0x12
	MINUS = 0x13
	DIV = 0x14
	ASSIGN = 0x15
	QUICK_ASSIGN = 0x16
	EQUALS = 0x17
	NOT_EQUALS = 0x18
	LESS_THAN = 0x19
	GREAT_THAN = 0x1A
	LESS_EQUAL = 0x1B
	GREAT_EQUAL = 0x1C
	RANGE = 0x1D        // ..
	VARIADIC = 0x1E		// ...


	IDENTIFIER = 0x1F
)

var (
	TokStrings = [...]string{"","{","}","[","]","(",")",",",".",";",":","?","*","!","?","&","|","~","+","-","/","=",":=","==","!=","<",">","<=",">=","..","...","identifier"}
)

type Token struct {
	Type uint32
	Value interface{}
}

func (token *Token) Debug() {
	Log(TokStrings[token.Type], token.Value)
}