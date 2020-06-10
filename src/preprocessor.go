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

	return preprocessor.Run()
}

func (preprocessor *Preprocessor) Run() string {
	preprocessor.Consumer.Scanner.Scan()
	for !preprocessor.Consumer.End() {
		r := preprocessor.Consumer.Peek()
		switch r {
		case rune('\n'):
			fallthrough
		case rune('\r'):
			preprocessor.Newline()
			break
		case rune('/'):
			preprocessor.Comment()
			break
		default:
			preprocessor.Compiler.Critical(preprocessor.Consumer.Reporter, ERR_UNEXPECTED_CHAR, "unexpected character")
			break
		}
	}

	return preprocessor.Consumer.Source
}

func (preprocessor *Preprocessor) Newline() {
	preprocessor.Consumer.Reporter.Position.Line++
	preprocessor.Consumer.Reporter.Position.Indent = 0
	preprocessor.Consumer.Advance()

	// get the next line of text
	preprocessor.Consumer.Scanner.Scan()
	preprocessor.Consumer.Reporter.CurrentLine = preprocessor.Consumer.Scanner.Text()
}

func (preprocessor *Preprocessor) Comment() {
	preprocessor.Consumer.Advance()
	r := preprocessor.Consumer.Peek()
	// single line comment
	if r == rune('/') {
		for r != rune('\n') && r != rune('\r') {
			r = preprocessor.Consumer.Advance()
		}
		preprocessor.Newline()
	} else if r == rune('*') {
		// multiline comment
	}
}
