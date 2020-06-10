package src

type LexConsumer struct {
	Source string
}

func (lexConsumer *LexConsumer) Peek() rune {
	return '0'
}

func (lexConsumer *LexConsumer) Consume(char rune) rune {
	return '0'
}

func (lexConsumer *LexConsumer) Advance() rune {
	return '0'
}
