package src

type Directives struct {
	Compiler  *Compiler
	Consumer  *ParseConsumer
	SymTable  *SymTable
	NewTokens []*Token
}

func ProcessDirectives(compiler *Compiler, tokens []*Token) []*Token {
	reporter := NewReporter(compiler.File.Filename, compiler.File.Source)
	consumer := NewParseConsumer(tokens, reporter, compiler)

	directives := Directives{
		Compiler: compiler,
		Consumer: consumer,
		SymTable: NewSymTable(nil),
	}

	result := directives.Run()

	return result
}

func (directives *Directives) Run() []*Token {
	for !directives.Consumer.End() {
		t := directives.Consumer.Peek().Type
		switch t {
		default:
			directives.Consumer.Advance()
			break
			//directives.Compiler.Critical(directives.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "token wasn't expected")
		}
	}

	return directives.Consumer.Tokens
}
