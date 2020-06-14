package src

const (
	ERR_UNEXPECTED_TOKEN = 0x0
)

// buffer the current active directives so we can process the rest of the tokens
type DirectiveBuf struct {
	Modifiers uint32
}

type Parser struct {
	Compiler     *Compiler
	Consumer     *ParseConsumer
	SymTable     *SymTable
	Root 		 *RootAST
	DirectiveBuf *DirectiveBuf
}

func Parse(compiler *Compiler, tokens []*Token) *RootAST {
	reporter := NewReporter(compiler.FileName, compiler.Source)
	consumer := NewParseConsumer(tokens, reporter, compiler)

	parser := Parser{
		Compiler: compiler,
		Consumer: consumer,
		SymTable: NewSymTable(),
	}
	parser.Root = parser.Run()
	return parser.Root
}

func (parser *Parser) Run() *RootAST {
	Root := &RootAST{}
	for !parser.Consumer.End() {
		t := parser.Consumer.Advance()
		// any top level expression is an identifier
		switch t.Type {
		case IDENTIFIER:
			Root.Statements = append(Root.Statements, parser.Define(t))
		default:
			parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "unexpected token")
		}
	}
	return Root
}

// parse a struct
func (parser *Parser) Struct(identifier *Token) AST {
	s := &StructAST{Identifier: identifier,}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' after 'struct'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	return s
}

// parse a function
func (parser *Parser) Fun(identifier *Token) AST {
	f := &FnAST{Identifier: identifier,RetType: *parser.ParseType()}

	if parser.Consumer.Consume(LEFT_PAREN) != nil {

		// process the arguments

		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
	}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' after ')'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	return f
}

// parse a variable definition (this only includes the identifier and type e.g. X : i32;, and assigning to
// a definition e.g. X : i32 = 1;)
func (parser *Parser) Define(identifier *Token) AST{
	def := &VarDefAST{
		Identifier: identifier,
		Type:       TavType{},
		Assignment: nil,
	}
	// explicit type define
	if parser.Consumer.Consume(COLON) != nil{
		// check if we have a pointer
		def.Type = *parser.ParseType()
		switch def.Type.Type {
		case STRUCT:
			return parser.Struct(identifier)
		case FN:
			return parser.Fun(identifier)
		default:
			if parser.Consumer.Consume(ASSIGN) != nil {
				def.Assignment = parser.Expression()
			}
		}
	}else{
		parser.Consumer.Consume(QUICK_ASSIGN)
		// we are doing a quick assign
		t, assignment := parser.QuickAssign()
		def.Type = t
		def.Assignment = assignment
	}
	parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after assignment")
	return def
}

// FOR NOW, WE DON'T SUPPORT QUICK ASSIGNING STRUCTS OR FUNCTIONS
// parse a variable quick assign (e.g. X := 1)
func (parser *Parser) QuickAssign() (TavType, AST){
	return TavType{
		Type: ANY,
		IsPtr:  false,
		PtrVal: nil,
	}, parser.Expression()
}

func (parser *Parser) Expression() AST {
	return nil
}

// figure out the type of an expression
// TODO For now this will return ANY type
func (parser *Parser) FigureType(expression ExprStmtAST) uint32{
	return ANY
}

// parse a type
func (parser *Parser) ParseType() *TavType{
	typ := &TavType{
		Type:   0,
		IsPtr:  false,
		PtrVal: nil,
	}
	// if it is a pointer, recursively get the pointer value
	if parser.Consumer.Consume(STAR) != nil{
		typ.IsPtr = true
		typ.PtrVal = parser.ParseType()
	}else{
		if t:=parser.Consumer.Consume(TYPE); t!=nil{
			// it isn't a pointer, so get the type
			typ.Type = t.Value.(uint32);
		}
	}
	return typ
}