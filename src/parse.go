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
	reporter := NewReporter(compiler.Source)
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
			Root.Statements = append(Root.Statements, parser.Identifier(t))
		default:
			parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "unexpected token")
		}
	}
	return Root
}

func (parser *Parser) Identifier(identifier *Token) AST {
	parser.Consumer.ConsumeErr(COLON, ERR_UNEXPECTED_TOKEN, "expected ':'")

	switch parser.Consumer.Advance().Type {
	case STRUCT:
		return parser.Struct(identifier)
	case FN:
		return parser.Fun(identifier)
	default:
		parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "expected token after identifier")
	}
	return nil
}

func (parser *Parser) Struct(identifier *Token) AST {
	s := &StructAST{Identifier: identifier}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}'")
	return s
}

func (parser *Parser) Fun(identifier *Token) AST {
	f := &FunAST{Identifier: identifier}
	// function with no arguments
	if parser.Consumer.Consume(LEFT_PAREN) != nil {
		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected ')'")
	}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}'")
	return f
}
