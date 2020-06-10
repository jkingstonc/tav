package src

type Preprocessor struct {
	Compiler *Compiler
	Consumer *LexConsumer
	SymTable *SymTable
}

func Preprocess(compiler *Compiler, source string) string {

	reporter := NewReporter()
	consumer := NewLexConsumer(source, reporter)
	symtable := NewSymTable()

	preprocessor := &Preprocessor{
		Compiler: compiler,
		Consumer: consumer,
		SymTable: symtable,
	}

	return source
}
