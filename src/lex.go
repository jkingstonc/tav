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
		r := lexer.Consumer.Advance()
		switch r {
		case '\n':
			fallthrough
		case '\r':
			lexer.Newline()
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
		case '}':
			lexer.Tok(RIGHT_CURLY, nil)
		case '[':
			lexer.Tok(LEFT_BRACKET, nil)
		case ']':
			lexer.Tok(RIGHT_BRACKET, nil)
		case '(':
			lexer.Tok(LEFT_PAREN, nil)
		case ')':
			lexer.Tok(RIGHT_PAREN, nil)
		case ',':
			lexer.Tok(COMMA, nil)
		case '.':
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
		case ':':
			if lexer.Consumer.Consume('=') {
				lexer.Tok(QUICK_ASSIGN, nil)
			} else {
				lexer.Tok(COLON, nil)
			}
		case '?':
			lexer.Tok(QUESTION, nil)
		case '*':
			lexer.Tok(STAR, nil)
		case '!':
			if lexer.Consumer.Consume('=') {
				lexer.Tok(NOT_EQUALS, nil)
			} else {
				lexer.Tok(BANG, nil)
			}
		case '<':
			if lexer.Consumer.Consume('=') {
				lexer.Tok(LESS_EQUAL, nil)
			} else {
				lexer.Tok(LESS_THAN, nil)
			}
		case '>':
			if lexer.Consumer.Consume('=') {
				lexer.Tok(GREAT_EQUAL, nil)
			} else {
				lexer.Tok(GREAT_THAN, nil)
			}
		case '%':
			lexer.Tok(PERCENT, nil)
		case '&':
			lexer.Tok(BIN_AND, nil)
		case '|':
			lexer.Tok(BIN_OR, nil)
		case '~':
			lexer.Tok(WIGGLE, nil)
		case '+':
			lexer.Tok(PLUS, nil)
		case '-':
			lexer.Tok(MINUS, nil)
		case '=':
			if lexer.Consumer.Consume('=') {
				lexer.Tok(EQUALS, nil)
			} else {
				lexer.Tok(ASSIGN, nil)
			}
		case '@':
			lexer.Tok(ADDR, nil)
		default:
			if !(IsChar(r) && lexer.Identifier(r)) {
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

func (lexer *Lexer) Identifier(r rune) bool {
	switch r {
	case 'u':
		if !lexer.CheckKeyword("8", U8) {
			if !lexer.CheckKeyword("16", U16) {
				if !lexer.CheckKeyword("32", U32) {
					if !lexer.CheckKeyword("64", U64) {
						return false
					}
				}
			}
		}
		return true
	case 'i':
		if !lexer.CheckKeyword("8", I8) {
			if !lexer.CheckKeyword("16", I16) {
				if !lexer.CheckKeyword("32", I32) {
					if !lexer.CheckKeyword("64", I64) {
						if !lexer.CheckKeyword("f", IF) {
							if !lexer.CheckKeyword("fdef", IFDEF) {
								return false
							}
						}
					}
				}
			}
		}
		return true
	case 'f':
		if !lexer.CheckKeyword("32", F32) {
			if !lexer.CheckKeyword("64", f64) {
				if !lexer.CheckKeyword("un", FUN) {
					if !lexer.CheckKeyword("or", FOR) {
						return false
					}
				}
			}
		}
		return true
	case 'b':
		if !lexer.CheckKeyword("ool", BOOL) {
			if !lexer.CheckKeyword("reak", BREAK) {
				return false
			}
		}
		return true
	case 's':
		if !lexer.CheckKeyword("tring", STRING) {
			if !lexer.CheckKeyword("truct", STRUCT) {
				if !lexer.CheckKeyword("witch", SWITCH) {
					return false
				}
			}
		}
		return true
	case 'a':
		if !lexer.CheckKeyword("ny", ANY) {
			if !lexer.CheckKeyword("nd", AND) {
				return false
			}
		}
		return true
	case 'n':
		if !lexer.CheckKeyword("ull", NULL) {
			return false
		}
		return true
	case 'd':
		if !lexer.CheckKeyword("ef", DEF) {
			return false
		}
	case 'r':
		if !lexer.CheckKeyword("un", RUN) {
			if !lexer.CheckKeyword("eturn", RETURN) {
				return false
			}
		}
		return true
	case 'e':
		if !lexer.CheckKeyword("lif", ELIF) {
			if !lexer.CheckKeyword("lse", ELSE) {
				if !lexer.CheckKeyword("xpose", EXPOSE) {
					return false
				}
			}
		}
		return true
	case 'c':
		if !lexer.CheckKeyword("ase", CASE) {
			if !lexer.CheckKeyword("ontinue", CONTINUE) {
				return false
			}
		}
		return true
	case 'h':
		if !lexer.CheckKeyword("ide", HIDE) {
			return false
		}
		return true
	case 'o':
		if !lexer.CheckKeyword("r", OR) {
			return false
		}
		return true
	}

	var identifier strings.Builder
	identifier.WriteRune(r)
	for !lexer.Consumer.End() && (IsChar(lexer.Consumer.Peek()) || IsNum(lexer.Consumer.Peek())) {
		identifier.WriteRune(lexer.Consumer.Peek())
		lexer.Consumer.Advance()
	}
	lexer.Tok(IDENTIFIER, identifier.String())
	return true
}

func (lexer *Lexer) CheckKeyword(keyword string, tok uint32) bool {
	l := len(keyword)
	// check if the length is greater than the source string
	if int(lexer.Consumer.Counter)+l > len(*lexer.Consumer.Source) {
		return false
	}
	valid := true
	// compare to see if the runes match in the rest of the string
	for i := 0; i < l; i++ {
		if keyword[i] != (*lexer.Consumer.Source)[int(lexer.Consumer.Counter)+i] {
			valid = false
		}
	}
	// if they match then advance and return true
	if valid {
		lexer.Consumer.AdvanceMul(uint32(l))
		lexer.Tok(tok, nil)
	}
	return valid
}
