package src

const (
	LEFT_CURLY    uint32 = 0x1
	RIGHT_CURLY   uint32 = 0x2
	LEFT_BRACKET  uint32 = 0x3
	RIGHT_BRACKET uint32 = 0x4
	LEFT_PAREN    uint32 = 0x5
	RIGHT_PAREN   uint32 = 0x6
	COMMA         uint32 = 0x7
	PERIOD        uint32 = 0x8
	SEMICOLON     uint32 = 0x9
	COLON         uint32 = 0xA
	QUESTION      uint32 = 0xB
	STAR          uint32 = 0xC
	BANG          uint32 = 0xD
	PERCENT       uint32 = 0xE
	BIN_AND       uint32 = 0xF
	BIN_OR        uint32 = 0x10
	WIGGLE        uint32 = 0x11
	PLUS          uint32 = 0x12
	MINUS         uint32 = 0x13
	DIV           uint32 = 0x14
	ASSIGN        uint32 = 0x15
	QUICK_ASSIGN  uint32 = 0x16
	EQUALS        uint32 = 0x17
	NOT_EQUALS    uint32 = 0x18
	LESS_THAN     uint32 = 0x19
	GREAT_THAN    uint32 = 0x1A
	LESS_EQUAL    uint32 = 0x1B
	GREAT_EQUAL   uint32 = 0x1C
	RANGE         uint32 = 0x1D // ..
	VARIADIC      uint32 = 0x1E // ...

	// variables
	IDENTIFIER uint32 = 0x1F
	ADDR       uint32 = 0x20 // @
	TYPE       uint32 = 0x21
	NULL       uint32 = 0x22
	TRUE       uint32 = 0x23
	FALSE      uint32 = 0x24
	SLITERAL   uint32 = 0x25
	NLITERAL   uint32 = 0x26

	// directives
	NATIVE uint32 = 0x27
	DEF    uint32 = 0x28
	RUN    uint32 = 0x29
	IFDEF  uint32 = 0x2A
	ENDIF  uint32 = 0x2B
	HIDE   uint32 = 0x2C
	EXPOSE uint32 = 0x2D
	PACK   uint32 = 0x2E
	IMPORT uint32 = 0x2F

	IF       uint32 = 0x30
	ELIF     uint32 = 0x31
	ELSE     uint32 = 0x32
	FOR      uint32 = 0x33
	SWITCH   uint32 = 0x34
	CASE     uint32 = 0x35
	BREAK    uint32 = 0x36
	CONTINUE uint32 = 0x37
	RETURN   uint32 = 0x38
	AND      uint32 = 0x39
	OR       uint32 = 0x3A
	SLEFT    uint32 = 0x3B
	SRIGHT   uint32 = 0x3C
)

var (
	TokStrings = [...]string{"", "{", "}", "[", "]", "(", ")", ",", ".", ";", ":", "?", "*", "!", "?",
		"&", "|", "~", "+", "-", "/", "=", ":=", "==", "!=", "<", ">", "<=", ">=", "..", "...", "identifier", "@", "type",
		"null","true","false", "sliteral", "nliteral", "native","def", "run", "ifdef", "endif", "hide", "expose", "pack",
		"import", "if", "elif", "else", "for", "switch", "case", "break", "continue", "return", "<<", ">>"}
)

type Token struct {
	Position Position
	Type  	 uint32
	Value 	 interface{}
}

func (token *Token) Debug() {
	Log(TokStrings[token.Type], token.Value)
}

func (token *Token) Lexme() string{
	return token.Value.(string)
}
