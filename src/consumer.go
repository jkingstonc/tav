package src

type LexConsumer struct {
	Source   string
	Reporter LexReporter
}

func NewLexConsumer(source string, reporter LexReporter) LexConsumer {
	return LexConsumer{
		Source:   source,
		Reporter: reporter,
	}
}

func (lexConsumer *LexConsumer) Peek() rune {
	return '0'
}

func (lexConsumer *LexConsumer) Consume(char rune) rune {
	return '0'
}

func (lexConsumer *LexConsumer) Advance() rune {
	lexConsumer.Reporter.Position.Indent++
	return '0'
}
