package src

import "time"

const (
	ERR_MULTIPLE_SEMICOLON = 0x0
	ERR_UNEXPECTED_CHAR    = 0x1
)

type Lexer struct {
	Compiler *Compiler
	Consumer LexConsumer
}

func Lex(compiler *Compiler, source string) []Token {
	start := time.Now()

	// lexer := &Lexer{
	// 	Compiler: compiler,
	// 	Consumer: LexConsumer{},
	// }

	// lexer.Reporter.CurrentLine = "i := 1;;"
	// lexer.Reporter.Position.Indent = 8
	// lexer.Reporter.Position.Line = 0
	// lexer.Compiler.Report(WARNING, lexer.Reporter, ERR_MULTIPLE_SEMICOLON, "multiple semi-colons")

	// lexer.Reporter.CurrentLine = "i := 1.23!45;"
	// lexer.Reporter.Position.Indent = 10
	// lexer.Reporter.Position.Line = 0
	// lexer.Compiler.Report(CRITICAL, lexer.Reporter, ERR_UNEXPECTED_CHAR, "unexpected character")

	end := time.Since(start)
	Log("front end took ", end.Seconds(), "ms")
	return nil
}
