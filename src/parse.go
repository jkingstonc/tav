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
	Root 		 AST
	DirectiveBuf *DirectiveBuf
}

func Parse(compiler *Compiler, tokens []*Token) AST {
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

func (parser *Parser) Run() AST {

	Root := &RootAST{}

	for !parser.Consumer.End() {
		t := parser.Consumer.Advance()
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
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	return s
}

// parse a function
func (parser *Parser) Fun(identifier *Token) AST {
	f := &FnAST{Identifier: identifier,}
	if t:=parser.Consumer.Consume(TYPE); t!=nil{
		f.RetType = t.Value.(uint32)
	}
	if parser.Consumer.Consume(LEFT_PAREN) != nil {

		// process the arguments

		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
	}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	return f
}

// parse a variable definition (this only includes the identifier and type e.g. X : i32;, and assigning to
// a definition e.g. X : i32 = 1;)
func (parser *Parser) Define(identifier *Token) AST{
	def := &DefineAST{
		Identifier: identifier,
		Type:       0,
		Assignment: nil,
	}
	// explicit type define
	if parser.Consumer.Consume(COLON) != nil{
		t := parser.Consumer.ConsumeErr(TYPE, ERR_UNEXPECTED_TOKEN, "expected type after ':'")
		def.Type = t.Value.(uint32)
		switch def.Type {
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
func (parser *Parser) QuickAssign() (uint32, AST){
	return ANY, parser.Expression()
}

func (parser *Parser) Expression() AST {
	return nil
}

// figure out the type of an expression
// TODO For now this will return ANY type
func (parser *Parser) FigureType(expression ExprAST) uint32{
	return ANY
}