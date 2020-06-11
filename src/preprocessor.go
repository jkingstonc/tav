package src

import "strings"

type Preprocessor struct {
	Compiler *Compiler
	Consumer *LexConsumer
	SymTable *SymTable
	NewStr   strings.Builder
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
			preprocessor.Consumer.Advance()

			// issue here is that the lexer then doesn't know there is a line here

			break
		case rune('/'):
			preprocessor.Comment()
			break
		default:
			preprocessor.NewStr.WriteRune(r)
			preprocessor.Consumer.Advance()
			break
		}
	}

	return preprocessor.NewStr.String()
}

func (preprocessor *Preprocessor) Newline() {
	preprocessor.Consumer.Reporter.Position.Line++
	preprocessor.Consumer.Reporter.Position.Indent = 0

	// get the next line of text
	preprocessor.Consumer.Scanner.Scan()
	preprocessor.Consumer.Reporter.CurrentLine = preprocessor.Consumer.Scanner.Text()
}

func (preprocessor *Preprocessor) Comment() {
	preprocessor.Consumer.Advance()
	//r := preprocessor.Consumer.Peek()
	// single line comment
	if preprocessor.Consumer.Expect('/') {
		for !preprocessor.Consumer.End() && !preprocessor.Consumer.Expect('\n') && !preprocessor.Consumer.Expect('\r') {
			preprocessor.Consumer.Advance()
		}
		preprocessor.Newline()
	} else if preprocessor.Consumer.Expect('*') {

	}
	// preprocessor.Consumer.Advance()
	// r := preprocessor.Consumer.Peek()
	// // single line comment
	// if r == rune('/') {
	// 	for !preprocessor.Consumer.End() && r != rune('\n') && r != rune('\r') {
	// 		r = preprocessor.Consumer.Advance()
	// 	}
	// 	preprocessor.Newline()
	// } else if r == rune('*') {
	// }
}
