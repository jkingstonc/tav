package src

import (
	"strings"
)

const (
	ERR_UNEXPECTED_CHAR 	   = 0x0
	ERR_STRING_ESCAPED  	   = 0x1
	ERR_INVALID_NUMBER_LITERAL = 0x2
)

type Lexer struct {
	Compiler *Compiler
	Consumer *LexConsumer
	Tokens   []*Token
}

func Lex(compiler *Compiler) []*Token {
	reporter := NewReporter(compiler.FileName, compiler.Source)
	consumer := NewLexConsumer(compiler.Source, reporter)

	lexer := Lexer{
		Compiler: compiler,
		Consumer: consumer,
		Tokens:   nil,
	}

	result := lexer.Run()
	return result
}

func (lexer *Lexer) Run() []*Token {
	for !lexer.Consumer.End() {
		lexer.Consumer.SkipWhitespace()
		r := lexer.Consumer.Advance()
		switch r {
		case '\n':  // newline
			lexer.Newline()
		case '\r':	// carrige return
			lexer.Consumer.Reporter.Position.Indent = 1
			// issue here is that the lexer then doesn't know there is a line here
			break
		case '/':
			if lexer.Consumer.Consume('/') {
				lexer.LineComment()
			} else if lexer.Consumer.Consume('*') {
				lexer.BlockComment()
			} else {
				lexer.Tok(DIV, nil)
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
				if IsNum(lexer.Consumer.Peek()){
					lexer.NumberLiteral(r)
				}else {
					lexer.Tok(PERIOD, nil)
				}
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
			}else if lexer.Consumer.Consume('<'){
				lexer.Tok(SLEFT, nil)
			} else {
				lexer.Tok(LESS_THAN, nil)
			}
		case '>':
			if lexer.Consumer.Consume('=') {
				lexer.Tok(GREAT_EQUAL, nil)
			} else if lexer.Consumer.Consume('>'){
				lexer.Tok(SRIGHT, nil)
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
		case '\'':
			fallthrough
		case '"':
			lexer.StringLiteral(r)
		case '#':
			if !lexer.CheckKeyword("def", DEF, nil) {
				if !lexer.CheckKeyword("run", RUN, nil) {
					if !lexer.CheckKeyword("ifdef", IFDEF, nil) {
						if !lexer.CheckKeyword("endif", ENDIF, nil) {
							if !lexer.CheckKeyword("hide", HIDE, nil) {
								if !lexer.CheckKeyword("pack", PACK, nil) {
									if !lexer.CheckKeyword("expose", EXPOSE, nil) {
										if !lexer.CheckKeyword("import", IMPORT, nil) {
											if !lexer.CheckKeyword("native", NATIVE, nil) {
												lexer.Compiler.Critical(lexer.Consumer.Reporter, ERR_UNEXPECTED_CHAR, "unexpected character")
											}
										}
									}
								}
							}
						}
					}
				}
			}
		default:
			if IsChar(r){
				lexer.Identifier(r)
			}else if IsNum(r){
				lexer.NumberLiteral(r)
			}else{
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
	for !lexer.Consumer.End(){
		if lexer.Consumer.Consume('*'){
			if lexer.Consumer.Consume('/'){
				return
			}
		}
		lexer.Consumer.Advance()
	}
}

func (lexer *Lexer) Tok(tok uint32, val interface{}) {
	t := &Token{lexer.Consumer.Reporter.Position,tok, val}
	lexer.Tokens = append(lexer.Tokens, t)
}

func (lexer *Lexer) StringLiteral(r rune){
	s := strings.Builder{}
	for !lexer.Consumer.End() && lexer.Consumer.Peek() != r{
		s.WriteRune(lexer.Consumer.Advance())
		if lexer.Consumer.End(){
			lexer.Compiler.Critical(lexer.Consumer.Reporter, ERR_STRING_ESCAPED, "string must be closed with ' or \"")
		}
	}
	lexer.Consumer.Consume(r)
	lexer.Tok(SLITERAL, s.String())
}

func (lexer *Lexer) NumberLiteral(r rune) bool{
	hadPeriod := false
	s := strings.Builder{}
	if r == '.'{
		hadPeriod = true
		s.WriteRune('0')
	}
	s.WriteRune(r)
	for !lexer.Consumer.End() && (IsNum(lexer.Consumer.Peek()) || lexer.Consumer.Expect('.')){
		n := lexer.Consumer.Advance()
		if n == '.' && hadPeriod == false{
			hadPeriod = true
		}else if n == '.' && hadPeriod == true{
			lexer.Compiler.Critical(lexer.Consumer.Reporter, ERR_INVALID_NUMBER_LITERAL, "number cannot have more than 1 '.'")
			return false
		}
		s.WriteRune(n)

	}
	lexer.Tok(NLITERAL, s.String())
	return true
}

func (lexer *Lexer) Identifier(r rune) bool {
	switch r {
	case 'u':
		if !lexer.CheckKeyword("8", TYPE, TYPE_U8) {
			if !lexer.CheckKeyword("16", TYPE, TYPE_U16) {
				if !lexer.CheckKeyword("32", TYPE, TYPE_U32) {
					if !lexer.CheckKeyword("64", TYPE, TYPE_U64) {
						break
					}
				}
			}
		}
		return true
	case 'i':
		if !lexer.CheckKeyword("8", TYPE, TYPE_I8) {
			if !lexer.CheckKeyword("16", TYPE, TYPE_I16) {
				if !lexer.CheckKeyword("32", TYPE, TYPE_I32) {
					if !lexer.CheckKeyword("64", TYPE, TYPE_I64) {
						if !lexer.CheckKeyword("f", IF, nil) {
							break
						}
					}
				}
			}
		}
		return true
	case 'f':
		if !lexer.CheckKeyword("32", TYPE, TYPE_F32) {
			if !lexer.CheckKeyword("64", TYPE, TYPE_F64) {
				if !lexer.CheckKeyword("n", TYPE, TYPE_FN) {
					if !lexer.CheckKeyword("or", FOR, nil) {
						if !lexer.CheckKeyword("alse", FALSE, false) {
							break
						}
					}
				}
			}
		}
		return true
	case 'b':
		if !lexer.CheckKeyword("ool", TYPE, TYPE_BOOL) {
			if !lexer.CheckKeyword("reak", BREAK, nil) {
				break
			}
		}
		return true
	case 's':
		if !lexer.CheckKeyword("tring", TYPE, TYPE_STRING) {
			if !lexer.CheckKeyword("truct", TYPE, TYPE_STRUCT) {
				if !lexer.CheckKeyword("witch", SWITCH, nil) {
					break
				}
			}
		}
		return true
	case 'a':
		if !lexer.CheckKeyword("ny", TYPE, TYPE_ANY) {
			if !lexer.CheckKeyword("nd", AND, nil) {
				break
			}
		}
		return true
	case 'n':
		if !lexer.CheckKeyword("ull", NULL, nil) {
			break
		}
		return true
	case 'r':
		if !lexer.CheckKeyword("eturn", RETURN, nil) {
			break
		}
		return true
	case 'e':
		if !lexer.CheckKeyword("lif", ELIF, nil) {
			if !lexer.CheckKeyword("lse", ELSE, nil) {
				break
			}
		}
		return true
	case 'c':
		if !lexer.CheckKeyword("ase", CASE, nil) {
			if !lexer.CheckKeyword("ontinue", CONTINUE, nil) {
				break
			}
		}
		return true
	case 'o':
		if !lexer.CheckKeyword("r", OR, nil) {
			break
		}
		return true
	case 'p':
		if !lexer.CheckKeyword("ack", PACK, nil) {
			break
		}
		return true
	case 't':
		if !lexer.CheckKeyword("rue", TRUE, false) {
			break
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

func (lexer *Lexer) CheckKeyword(keyword string, tok uint32, value interface{}) bool {
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
		lexer.Tok(tok, value)
	}
	return valid
}
