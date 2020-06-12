package src

type Consumer struct {
	Reporter *Reporter
	// used to indicate which char/token we are at
	Counter uint32
}

type LexConsumer struct {
	Consumer
	Source  *string
}

type ParseConsumer struct {
	Consumer
	Tokens []*Token
}

func NewLexConsumer(source *string, reporter *Reporter) *LexConsumer {
	return &LexConsumer{
		Consumer: Consumer{
			Reporter: reporter,
		},
		Source: source,
	}
}

func NewParseConsumer(tokens []*Token, reporter *Reporter) *ParseConsumer {
	return &ParseConsumer{
		Consumer: Consumer{
			Reporter: reporter,
		},
		Tokens: tokens,
	}
}

func (lexConsumer *LexConsumer) SkipWhitespace() {
	for lexConsumer.Expect(' ') {
		lexConsumer.Advance()
	}
}

func (lexConsumer *LexConsumer) Peek() rune {
	return rune((*lexConsumer.Source)[lexConsumer.Counter])
}

func (lexConsumer *LexConsumer) Expect(r rune) bool {
	if !lexConsumer.End() {
		return lexConsumer.Peek() == r
	}
	return false
}

func (lexConsumer *LexConsumer) Consume(char rune) bool {
	if lexConsumer.Expect(char){
		lexConsumer.Advance()
		return true
	}
	return false
}

func (lexConsumer *LexConsumer) Advance() rune {
	r := (*lexConsumer.Source)[lexConsumer.Counter]
	lexConsumer.Counter++
	lexConsumer.Reporter.Position.Indent++
	return rune(r)
}

func (lexConsumer *LexConsumer) AdvanceMul(ammount uint32) rune {
	r := (*lexConsumer.Source)[lexConsumer.Counter]
	lexConsumer.Counter += ammount
	lexConsumer.Reporter.Position.Indent += ammount
	return rune(r)
}

func (lexConsumer *LexConsumer) End() bool {
	return !(lexConsumer.Counter < uint32(len(*lexConsumer.Source)))
}

func (parseConsumer *ParseConsumer) Peek() *Token {
	return parseConsumer.Tokens[parseConsumer.Counter]
}

func (parseConsumer *ParseConsumer) Expect(tokenType uint32) bool {
	if !parseConsumer.End() {
		return parseConsumer.Peek().Type == tokenType
	}
	return false
}

func (parseConsumer *ParseConsumer) Consume(tokenType uint32) *Token {
	t := parseConsumer.Advance()
	return t
}

func (parseConsumer *ParseConsumer) Advance() *Token {
	t := parseConsumer.Tokens[parseConsumer.Counter]
	parseConsumer.Counter++
	parseConsumer.Reporter.Position.Indent++
	return t
}

func (parseConsumer *ParseConsumer) AdvanceMul(ammount uint32) *Token {
	t := parseConsumer.Tokens[parseConsumer.Counter]
	parseConsumer.Counter += ammount
	parseConsumer.Reporter.Position.Indent += ammount
	return t
}
func (parseConsumer *ParseConsumer) End() bool {
	return !(parseConsumer.Counter < uint32(len(parseConsumer.Tokens)))
}
func (parseConsumer *ParseConsumer) Remove() {
	parseConsumer.Tokens = append(parseConsumer.Tokens[:parseConsumer.Counter], parseConsumer.Tokens[parseConsumer.Counter+1:]...)
}
