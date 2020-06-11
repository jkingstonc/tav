package src

const (
	ERR_UNEXPECTED_TOKEN = 0x0
)

type Parser struct {
	Reporter *Reporter
	SymTable *SymTable
}

func Parse(compiler *Compiler, tokens []*Token) *AST {
	parser := Parser{
		Reporter: NewReporter(compiler.Source),
		SymTable: NewSymTable(),
	}
	r := parser.Run()
	return r
}

func (parser *Parser) Run() *AST {
	return nil
}
