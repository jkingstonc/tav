package src

type Consumer struct {
	Reporter *Reporter
}

type LexConsumer struct {
	Consumer
	Source string
}

func NewLexConsumer(source string, reporter *Reporter) *LexConsumer {

	reporter.CurrentLine = source

	return &LexConsumer{
		Consumer: Consumer{
			Reporter: reporter,
		},
		Source: source,
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
