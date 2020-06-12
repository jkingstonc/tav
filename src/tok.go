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

	// variables
	IDENTIFIER = 0x1F
	PTR		   = 0x20
	U8		   = 0x21
	I8		   = 0x22
	U16		   = 0x23
	I16		   = 0x24
	U32		   = 0x25
	I32		   = 0x26
	U64		   = 0x27
	I64		   = 0x28
	f64		   = 0x29
	BOOL	   = 0x2A
	STRING	   = 0x2B
	STRUCT     = 0xAC
	FUN        = 0xAD
	ANY        = 0xAE
	SLITERAL   = 0xAF
	NLITERAL   = 0xB0
	NULL       = 0xB1

	// directives
	DEF 	   = 0xB2
	RUN        = 0xB3
	IFDEF      = 0xB4
	HIDE       = 0xB5
	EXPOSE     = 0xB6

	IF		   = 0xB7
	ELIF	   = 0xB8
	ELSE	   = 0xB9
	DO	       = 0xBA
	FOR	       = 0xBB
	SWITCH	   = 0xBC
	BREAK	   = 0xBD
	CONTINUE   = 0xBE
	RETURN	   = 0xBF
)

var (
	TokStrings = [...]string{"","{","}","[","]","(",")",",",".",";",":","?","*","!","?",
		"&","|","~","+","-","/","=",":=","==","!=","<",">","<=",">=","..","...","identifier","define"}
)

type Token struct {
	Type uint32
	Value interface{}
}

func (token *Token) Debug() {
	Log(TokStrings[token.Type], token.Value)
}