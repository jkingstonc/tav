package src

const (
	ERR_UNEXPECTED_TOKEN = 0x0
)

type Parser struct {
	Compiler *Compiler
	Consumer *ParseConsumer
	SymTable *SymTable
	AST      *AST
}

func Parse(compiler *Compiler, tokens []*Token) *AST {
	reporter := NewReporter(compiler.Source)
	consumer := NewParseConsumer(tokens, reporter, compiler)

	parser := Parser{
		Compiler: compiler,
		Consumer: consumer,
		SymTable: NewSymTable(),
	}
	r := parser.Run()
	return r
}

func (parser *Parser) Run() *AST {
	for !parser.Consumer.End(){
		t := parser.Consumer.Advance()
		switch t.Type {
		case IDENTIFIER:
			parser.Identifier(t)
		default:
			parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "unexpected token")
		}
	}
	return nil
}

func (parser *Parser) Identifier(identifier *Token) *AST{
	parser.Consumer.ConsumeErr(COLON, ERR_UNEXPECTED_TOKEN, "expected ':'")

	switch parser.Consumer.Advance().Type{
		case STRUCT:
			return parser.Struct(identifier)
		case FUN:
			return parser.Fun(identifier)
		default:
			parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "expected token after identifier")
	}
	return nil
}

func (parser *Parser) Struct(identifier *Token) *AST{
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{'")

	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}'")
	return nil
}

func (parser *Parser) Fun(identifier *Token)*AST{
	// function with no arguments
	if parser.Consumer.Consume(LEFT_PAREN) != nil{
		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected ')'")
	}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{'")
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}'")
	return nil
}
