package src

type Preprocessor struct {
	Compiler *Compiler
	Reporter LexReporter
	Consumer LexConsumer
	SymTable *SymTable
}

func Preprocess(compiler *Compiler, source string) string {

	reporter := NewLexReporter()
	consumer := NewLexConsumer(source, reporter)

	preprocessor := &Preprocessor{
		Compiler: compiler,
		Reporter: LexReporter{},
		Consumer: consumer,
		SymTable: &SymTable{
			SymbolID: make(map[string]uint32),
			Symbols:  make(map[uint32]Symbol),
		},
	}

	preprocessor.Consumer.Advance()

	return source
}
