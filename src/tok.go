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
	U8         uint32 = 0x21
	I8         uint32 = 0x22
	U16        uint32 = 0x23
	I16        uint32 = 0x24
	U32        uint32 = 0x25
	I32        uint32 = 0x26
	F32        uint32 = 0x27
	U64        uint32 = 0x28
	I64        uint32 = 0x29
	F64        uint32 = 0x2A
	BOOL       uint32 = 0x2B
	STRING     uint32 = 0x2C
	STRUCT     uint32 = 0x2D
	FN         uint32 = 0x2E
	ANY        uint32 = 0x2F
	NULL       uint32 = 0x30

	SLITERAL   uint32 = 0x31
	NLITERAL   uint32 = 0x32

	// directives
	NATIVE uint32 = 0x33
	DEF    uint32 = 0x34
	RUN    uint32 = 0x35
	IFDEF  uint32 = 0x36
	ENDIF  uint32 = 0x37
	HIDE   uint32 = 0x38
	EXPOSE uint32 = 0x39
	PACK   uint32 = 0x3A
	IMPORT uint32 = 0x3B

	IF       uint32 = 0x3C
	ELIF     uint32 = 0x3D
	ELSE     uint32 = 0x3E
	FOR      uint32 = 0x3F
	SWITCH   uint32 = 0x40
	CASE     uint32 = 0x41
	BREAK    uint32 = 0x42
	CONTINUE uint32 = 0x43
	RETURN   uint32 = 0x44
	AND      uint32 = 0x45
	OR       uint32 = 0x46
)

var (
	TokStrings = [...]string{"", "{", "}", "[", "]", "(", ")", ",", ".", ";", ":", "?", "*", "!", "?",
		"&", "|", "~", "+", "-", "/", "=", ":=", "==", "!=", "<", ">", "<=", ">=", "..", "...", "identifier", "@", "type",
		"u8", "i8", "u16", "i16", "u32", "i32", "f32", "u64", "i64", "f64", "bool", "string", "struct", "fn",
		"any", "sliteral", "nliteral", "null", "native","def", "run", "ifdef", "endif", "hide", "expose", "pack",
		"import", "if", "elif", "else", "for", "switch", "case", "break", "continue", "return"}
)

type Token struct {
	Type  uint32
	Value interface{}
}

func (token *Token) Debug() {
	Log(TokStrings[token.Type], token.Value)
}
