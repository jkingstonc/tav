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
	Log("front end took ", end.Seconds(), "seconds")
	return result
}

func (lexer *Lexer) Run() []*Token {
	for !lexer.Consumer.End() {
		lexer.Consumer.SkipWhitespace()
		r := lexer.Consumer.Peek()
		switch r {
		case '\n':
			fallthrough
		case '\r':
			lexer.Newline()
			lexer.Consumer.Advance()
			// issue here is that the lexer then doesn't know there is a line here
			break
		case '/':
			lexer.Consumer.Advance()
			if lexer.Consumer.Consume('/') {
				lexer.LineComment()
			} else if lexer.Consumer.Consume('*') {
				lexer.BlockComment()
			} else {
				lexer.Tok(DIV, nil)
				lexer.Consumer.Advance()
			}
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
			if lexer.Consumer.Consume('.') {
				if lexer.Consumer.Consume('.') {
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
			if lexer.Consumer.Consume('=') {
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
			lexer.Consumer.Advance()
			if lexer.Consumer.Consume('=') {
				lexer.Tok(NOT_EQUALS, nil)
			} else {
				lexer.Tok(BANG, nil)
			}
		case '<':
			lexer.Consumer.Advance()
			if lexer.Consumer.Consume('=') {
				lexer.Tok(LESS_EQUAL, nil)
			} else {
				lexer.Tok(LESS_THAN, nil)
			}
		case '>':
			lexer.Consumer.Advance()
			if lexer.Consumer.Consume('=') {
				lexer.Tok(GREAT_EQUAL, nil)
			} else {
				lexer.Tok(GREAT_THAN, nil)
			}
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
				lexer.Tok(IDENTIFIER, lexer.Identifier())
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
	for !lexer.Consumer.End() && !lexer.Consumer.Expect('\n') && !lexer.Consumer.Expect('\r') {
		lexer.Consumer.Advance()
	}
}

func (lexer *Lexer) BlockComment() {
}

func (lexer *Lexer) Tok(tok uint32, val interface{}) {
	t := &Token{tok, val}
	t.Debug()
	lexer.Tokens = append(lexer.Tokens, t)
}

func (lexer *Lexer) Identifier() string {
	var identifier strings.Builder
	for !lexer.Consumer.End() && (IsChar(lexer.Consumer.Peek()) || IsNum(lexer.Consumer.Peek())) {
		identifier.WriteRune(lexer.Consumer.Peek())
		lexer.Consumer.Advance()
	}
	return identifier.String()
}
