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
	Counter uint32
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

func (lexConsumer *LexConsumer) Peek() rune {
	return rune(lexConsumer.Source[lexConsumer.Counter])
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

func (lexConsumer *LexConsumer) End() bool {
	return !(lexConsumer.Consumer.Counter < uint32(len(lexConsumer.Source)))
}
