package src

import "time"

const (
	ERR_MULTIPLE_SEMICOLON = 0x0
	ERR_UNEXPECTED_CHAR    = 0x1
)

type Lexer struct {
	Compiler *Compiler
	Consumer *LexConsumer
	Tokens   []*Token
}

func Lex(compiler *Compiler, source string) []*Token {
	start := time.Now()


	reporter := NewReporter()
	consumer := NewLexConsumer(source, reporter)

	lexer := Lexer{
		Compiler : compiler,
		Consumer : consumer,
		Tokens   : nil,
	}

	result := lexer.Run()
	end := time.Since(start)
	Log("front end took ", end.Seconds(), "ms")
	return result
}


func (lexer *Lexer) Run() []*Token {
	lexer.Consumer.Scanner.Scan()
	for !lexer.Consumer.End() {
		r := lexer.Consumer.Peek()
		lexer.Consumer.SkipWhitespace(r)
		switch r {
		case rune('{'): 
			lexer.Tok(LEFT_CURLY, nil)
			lexer.Consumer.Advance()
		case rune('}'): 
			lexer.Tok(RIGHT_CURLY, nil)
			lexer.Consumer.Advance()
		case rune('['): 
			lexer.Tok(LEFT_BRACKET, nil)
			lexer.Consumer.Advance()
		case rune(']'): 
			lexer.Tok(RIGHT_BRACKET, nil)
			lexer.Consumer.Advance()
		case rune('('): 
			lexer.Tok(LEFT_PAREN, nil)
			lexer.Consumer.Advance()
		case rune(')'): 
			lexer.Tok(RIGHT_PAREN, nil)
			lexer.Consumer.Advance()
		case rune(','): 
			lexer.Tok(COMMA, nil)
			lexer.Consumer.Advance()
		case rune('.'): 
		    lexer.Consumer.Advance()
			if lexer.Consumer.Expect('.'){
				lexer.Consumer.Advance()
				if lexer.Consumer.Expect('.'){
					lexer.Consumer.Advance()
					lexer.Tok(VARIADIC, nil)
				}else{
					lexer.Tok(RANGE, nil)
				}
			}else{
				lexer.Tok(PERIOD, nil)
			}
		case rune(';'): 
			lexer.Tok(SEMICOLON, nil)
			lexer.Consumer.Advance()
		case rune(':'): 
			lexer.Consumer.Advance()
			if lexer.Consumer.Expect('='){
				lexer.Consumer.Advance()
				lexer.Tok(QUICK_ASSIGN, nil)
			}else{
				lexer.Tok(COLON, nil)
			}
		case rune('?'): 
			lexer.Tok(QUESTION, nil)
			lexer.Consumer.Advance()
		case rune('*'): 
			lexer.Tok(STAR, nil)
			lexer.Consumer.Advance()
		case rune('!'): 
			lexer.Tok(BANG, nil)
			lexer.Consumer.Advance()
		case rune('%'): 
			lexer.Tok(PERCENT, nil)
			lexer.Consumer.Advance()
		case rune('&'): 
			lexer.Tok(BIN_AND, nil)
			lexer.Consumer.Advance()
		case rune('|'): 
			lexer.Tok(BIN_OR, nil)
			lexer.Consumer.Advance()
		case rune('~'): 
			lexer.Tok(WIGGLE, nil)
			lexer.Consumer.Advance()
		case rune('+'): 
			lexer.Tok(PLUS, nil)
			lexer.Consumer.Advance()
		case rune('-'): 
			lexer.Tok(MINUS, nil)
			lexer.Consumer.Advance()
		case rune('/'): 
			lexer.Tok(BANG, nil)
			lexer.Consumer.Advance()
		case rune('='): 
			lexer.Consumer.Advance()
			if lexer.Consumer.Expect('='){
				lexer.Tok(EQUALS, nil)
				lexer.Consumer.Advance()
			}else{
				lexer.Tok(ASSIGN, nil)
				lexer.Consumer.Advance()
			}
		default:
			if IsChar(r){
				Log("char!")
				lexer.Consumer.Advance()
			}else{
				Log("ufkc", string(r))
				lexer.Compiler.Critical(lexer.Consumer.Reporter, ERR_UNEXPECTED_CHAR, "unexpected character")
			}
		}
	}

	return lexer.Tokens
}

func (lexer *Lexer) Tok(tok uint32, val interface{}){
	t := &Token{tok, val}
	t.Debug()
	lexer.Tokens = append(lexer.Tokens, t)
}