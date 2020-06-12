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
		case STRUCT:
			parser.Struct()
		default:
			parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "unexpected token")
		}
	}
	return nil
}

func (parser *Parser) Struct(){
	parser.Consumer.ConsumeErr(LEFT_BRACKET, ERR_UNEXPECTED_TOKEN, "expected {")

	parser.Consumer.ConsumeErr(LEFT_BRACKET, ERR_UNEXPECTED_TOKEN, "expected }")
}
