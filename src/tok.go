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
	ADDR	   = 0x20   // &
	U8		   = 0x21
	I8		   = 0x22
	U16		   = 0x23
	I16		   = 0x24
	U32		   = 0x25
	I32		   = 0x26
	F32		   = 0x27
	U64		   = 0x28
	I64		   = 0x29
	f64		   = 0x2A
	BOOL	   = 0x2B
	STRING	   = 0x2C
	STRUCT     = 0x2D
	FUN        = 0x2E
	ANY        = 0x2F
	SLITERAL   = 0x30
	NLITERAL   = 0x31
	NULL       = 0x32

	// directives
	DEF 	   = 0x33
	RUN        = 0x34
	IFDEF      = 0x35
	HIDE       = 0x36
	EXPOSE     = 0x37

	IF		   = 0x38
	ELIF	   = 0x39
	ELSE	   = 0x3A
	FOR	       = 0x3C
	SWITCH	   = 0x3D
	CASE       = 0x3E
	BREAK	   = 0x3F
	CONTINUE   = 0x40
	RETURN	   = 0x41
	AND		   = 0x42
	OR		   = 0x43
)

var (
	TokStrings = [...]string{"","{","}","[","]","(",")",",",".",";",":","?","*","!","?",
		"&","|","~","+","-","/","=",":=","==","!=","<",">","<=",">=","..","...","identifier","&",
		"u8","i8","u16","i16","u32","i32","f32","u64","i64","f64","bool","string","struct","fun",
	"any","sliteral","nliteral","null","def","run","ifdef","endif","hide","expose","if","elif","else",
	"for","switch","case","break","continue","return"}
)

type Token struct {
	Type uint32
	Value interface{}
}

func (token *Token) Debug() {
	Log(TokStrings[token.Type], token.Value)
}