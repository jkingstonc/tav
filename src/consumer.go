package src

import (
	"bufio"
	"strings"
)

type Consumer struct {
	Reporter *Reporter
	// used to indicate which char/token we are at
	Counter uint32
}

type LexConsumer struct {
	Consumer
	Source  string
	Scanner *bufio.Scanner
}

func NewLexConsumer(source string, reporter *Reporter) *LexConsumer {

	reporter.CurrentLine = source

	return &LexConsumer{
		Consumer: Consumer{
			Reporter: reporter,
		},
		Source:  source,
		Scanner: bufio.NewScanner(strings.NewReader(source)),
	}
}

func (lexConsumer *LexConsumer) SkipWhitespace(r rune){
	for r == rune(' '){
		Log("skippp", rune(r))
		r = lexConsumer.Advance()
	}
}

func (lexConsumer *LexConsumer) Peek() rune {
	return rune(lexConsumer.Source[lexConsumer.Counter])
}

func (lexConsumer *LexConsumer) Expect(r rune) bool {
	return lexConsumer.Peek() == r
}

func (lexConsumer *LexConsumer) Consume(char rune) rune {
	r := lexConsumer.Advance()
	return r
}

func (lexConsumer *LexConsumer) Advance() rune {
	r := lexConsumer.Source[lexConsumer.Counter]
	lexConsumer.Counter++
	lexConsumer.Reporter.Position.Indent++
	return rune(r)
}

func (lexConsumer *LexConsumer) AdvanceMul(ammount uint32) rune {
	r := lexConsumer.Source[lexConsumer.Counter]
	lexConsumer.Counter+=ammount
	lexConsumer.Reporter.Position.Indent+=ammount
	return rune(r)
}

func (lexConsumer *LexConsumer) End() bool {
	return !(lexConsumer.Counter < uint32(len(lexConsumer.Source)))
}
