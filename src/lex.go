package src

import (
	"strings"
	"time"
)

const (
	ERR_MULTIPLE_SEMICOLON = 0x0
	ERR_UNEXPECTED_CHAR    = 0x1
)

type Lexer struct {
	Compiler *Compiler
	Consumer *LexConsumer
	Tokens   []*Token
}

func Lex(compiler *Compiler) []*Token {
	start := time.Now()

	reporter := NewReporter(compiler.Source)
	consumer := NewLexConsumer(compiler.Source, reporter)

	lexer := Lexer{
		Compiler: compiler,
		Consumer: consumer,
		Tokens:   nil,
	}

	result := lexer.Run()
	end := time.Since(start)
	Log("front end took ", end.Seconds(), "ms")
	return result
}

func (lexer *Lexer) Run() []*Token {
	for !lexer.Consumer.End() {
		r := lexer.Consumer.Peek()
		lexer.Consumer.SkipWhitespace(r)
		switch r {
		case '\n':
			fallthrough
		case '\r':
			lexer.Newline()
			lexer.Consumer.Advance()
			// issue here is that the lexer then doesn't know there is a line here
			break
		case '/':
			if lexer.Consumer.Expect('/') {
				lexer.LineComment()
			} else if lexer.Consumer.Expect('*') {
				lexer.BlockComment()
			}
			lexer.Tok(BANG, nil)
			lexer.Consumer.Advance()
			break
		case '{':
			lexer.Tok(LEFT_CURLY, nil)
			lexer.Consumer.Advance()
		case '}':
			lexer.Tok(RIGHT_CURLY, nil)
			lexer.Consumer.Advance()
		case '[':
			lexer.Tok(LEFT_BRACKET, nil)
			lexer.Consumer.Advance()
		case ']':
			lexer.Tok(RIGHT_BRACKET, nil)
			lexer.Consumer.Advance()
		case '(':
			lexer.Tok(LEFT_PAREN, nil)
			lexer.Consumer.Advance()
		case ')':
			lexer.Tok(RIGHT_PAREN, nil)
			lexer.Consumer.Advance()
		case ',':
			lexer.Tok(COMMA, nil)
			lexer.Consumer.Advance()
		case '.':
			lexer.Consumer.Advance()
			if lexer.Consumer.Expect('.') {
				lexer.Consumer.Advance()
				if lexer.Consumer.Expect('.') {
					lexer.Consumer.Advance()
					lexer.Tok(VARIADIC, nil)
				} else {
					lexer.Tok(RANGE, nil)
				}
			} else {
				lexer.Tok(PERIOD, nil)
			}
		case ';':
			lexer.Tok(SEMICOLON, nil)
			lexer.Consumer.Advance()
		case ':':
			lexer.Consumer.Advance()
			if lexer.Consumer.Expect('=') {
				lexer.Consumer.Advance()
				lexer.Tok(QUICK_ASSIGN, nil)
			} else {
				lexer.Tok(COLON, nil)
			}
		case '?':
			lexer.Tok(QUESTION, nil)
			lexer.Consumer.Advance()
		case '*':
			lexer.Tok(STAR, nil)
			lexer.Consumer.Advance()
		case '!':
			lexer.Tok(BANG, nil)
			lexer.Consumer.Advance()
		case '%':
			lexer.Tok(PERCENT, nil)
			lexer.Consumer.Advance()
		case '&':
			lexer.Tok(BIN_AND, nil)
			lexer.Consumer.Advance()
		case '|':
			lexer.Tok(BIN_OR, nil)
			lexer.Consumer.Advance()
		case '~':
			lexer.Tok(WIGGLE, nil)
			lexer.Consumer.Advance()
		case '+':
			lexer.Tok(PLUS, nil)
			lexer.Consumer.Advance()
		case '-':
			lexer.Tok(MINUS, nil)
			lexer.Consumer.Advance()
		case '=':
			lexer.Consumer.Advance()
			if lexer.Consumer.Expect('=') {
				lexer.Tok(EQUALS, nil)
				lexer.Consumer.Advance()
			} else {
				lexer.Tok(ASSIGN, nil)
			}
		default:
			if IsChar(r) {
				lexer.Tok(IDENTIFIER, lexer.Identifier(r))
			} else {
				lexer.Compiler.Critical(lexer.Consumer.Reporter, ERR_UNEXPECTED_CHAR, "unexpected character")
			}
		}
	}

	return lexer.Tokens
}

func (lexer *Lexer) Newline() {
	lexer.Consumer.Reporter.Position.Line++
	lexer.Consumer.Reporter.Position.Indent = 0
}

func (lexer *Lexer) LineComment() {
	lexer.Consumer.Advance()
	for !lexer.Consumer.End() && !lexer.Consumer.Expect('\n') && !lexer.Consumer.Expect('\r') {
		lexer.Consumer.Advance()
	}
	lexer.Newline()
}

func (lexer *Lexer) BlockComment() {
}

func (lexer *Lexer) Tok(tok uint32, val interface{}) {
	t := &Token{tok, val}
	t.Debug()
	lexer.Tokens = append(lexer.Tokens, t)
}

func (lexer *Lexer) Identifier(r rune) string {
	var identifier strings.Builder
	lexer.Consumer.Advance()
	identifier.WriteRune(r)
	r = lexer.Consumer.Peek()
	for IsChar(r) || IsNum(r) {
		identifier.WriteRune(r)
		lexer.Consumer.Advance()
		r = lexer.Consumer.Peek()
	}
	return identifier.String()
}
