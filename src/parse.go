package src

import (
	"strconv"
	"strings"
)

const (
	ERR_UNEXPECTED_TOKEN = 0x0
	ERR_INVALID_TYPE     = 0x1
)

// buffer the current active directives so we can process the rest of the tokens
type DirectiveBuf struct {
	Modifiers uint32
}

type Parser struct {
	Compiler     *Compiler
	Consumer     *ParseConsumer
	SymTable     *SymTable
	Root         *RootAST
	DirectiveBuf *DirectiveBuf
}

func Parse(compiler *Compiler, tokens []*Token) *RootAST {
	reporter := NewReporter(compiler.File.Filename, compiler.File.Source)
	consumer := NewParseConsumer(tokens, reporter, compiler)

	parser := Parser{
		Compiler: compiler,
		Consumer: consumer,
		SymTable: NewSymTable(),
	}
	parser.Root = parser.Run()
	return parser.Root
}

func (parser *Parser) Run() *RootAST {
	Root := &RootAST{}
	for !parser.Consumer.End() {
		t := parser.Consumer.Peek()
		// any top level expression is an identifier
		switch t.Type {
		case IDENTIFIER:
			Root.Statements = append(Root.Statements, parser.Define())
		default:
			parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "unexpected token")
		}
	}
	return Root
}

func (parser *Parser) Prelim() AST {
	return parser.Statement()
}

func (parser *Parser) Statement() AST {
	var ast AST
	if parser.Consumer.Expect(IDENTIFIER) && parser.Consumer.ExpectAhead(COLON,1) {
		ast = parser.Define()
	} else if parser.Consumer.Expect(IDENTIFIER) && parser.Consumer.ExpectAhead(QUICK_ASSIGN,1) {
		ast = parser.QuickAssign()
	} else if parser.Consumer.Expect(IDENTIFIER) {
		ast = parser.Assignment()
	} else if parser.Consumer.Consume(RETURN) != nil {
		ast = parser.Return()
	} else if parser.Consumer.Consume(BREAK) != nil {
		ast = parser.Break()
	} else if parser.Consumer.Consume(FOR) != nil {
		return parser.For()
	} else if parser.Consumer.Consume(IF) != nil {
		return parser.If()
	} else if parser.Consumer.Expect(LEFT_CURLY) {
		return &BlockAST{Statements: parser.ParseStmtBlock()}
	} else {
		ast = parser.ExpressionStmt()
	}
	parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' at end of statement")
	return ast
}

func (parser *Parser) ExpressionStmt() AST {
	return &ExprStmtAST{Expression: parser.Expression()}
}

func (parser *Parser) Return() AST {
	r := &ReturnAST{Value: parser.Expression()}
	return r
}

func (parser *Parser) Break() AST {
	r := &BreakAST{}
	return r
}

func (parser *Parser) For() AST {
	r := &ReturnAST{Value: parser.Expression()}
	return r
}

func (parser *Parser) If() AST {
	ifStmt := &IfAST{
		IfCondition:   nil,
		IfBody:        nil,
		ElifCondition: nil,
		ElifBody:      nil,
		ElseBody:      nil,
	}
	ifStmt.IfCondition = parser.Expression()
	ifStmt.IfBody = parser.Statement()
	for parser.Consumer.Consume(ELIF)!=nil{
		ifStmt.ElifCondition = append(ifStmt.ElifCondition,parser.Expression())
		ifStmt.ElifBody = append(ifStmt.ElifBody,parser.Statement())
	}
	if parser.Consumer.Consume(ELSE)!=nil{
		ifStmt.ElseBody = parser.Statement()
	}
	return ifStmt
}

func (parser *Parser) ParseStmtBlock() []AST {
	parser.SymTable.NewScope("block_body")
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' at start of statement block")
	var statements []AST
	for !parser.Consumer.Expect(RIGHT_CURLY) && !parser.Consumer.End() {
		statements = append(statements, parser.Statement())
	}
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}' at end of statement block")
	parser.SymTable.PopScope()
	return statements
}

// parse a struct
func (parser *Parser) Struct(identifier *Token) AST {
	name := identifier.Lexme()
	parser.SymTable.NewScope(name + "_members")

	s := &StructAST{Identifier: identifier}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' after 'struct'")

	for !parser.Consumer.Expect(RIGHT_CURLY) {
		member := parser.Define().(*VarDefAST)
		s.Fields = append(s.Fields, member)
		parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after member decleration")
	}

	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	parser.SymTable.PopScope()

	// add the identifier to the current symbol table
	parser.SymTable.Add(name, NewTavType(TYPE_STRUCT, "", 0, nil), nil)

	return s
}

// parse a function
func (parser *Parser) Fn(identifier *Token) AST { // add the identifier to the current symbol table

	name := identifier.Lexme()
	f := &FnAST{Identifier: identifier, RetType: *parser.ParseType()}

	parser.SymTable.Add(name, NewTavType(TYPE_FN, "", 0, &f.RetType), nil)

	parser.SymTable.NewScope(name + "_body")
	if parser.Consumer.Consume(LEFT_PAREN) != nil {
		var params []VarDefAST
		// process the arguments
		for !parser.Consumer.Expect(RIGHT_PAREN) {
			// each paramater is essentially a variable decleration
			params = append(params, *parser.Define().(*VarDefAST))
			if parser.Consumer.Expect(RIGHT_PAREN) {
				break
			}
			parser.Consumer.ConsumeErr(COMMA, ERR_UNEXPECTED_TOKEN, "expected ',' between parameters")
		}
		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
		f.Params = params
	}
	var statements []AST
	// parse the function body
	// this is not a statement block, we need the paramaters and the body in the name scope
	if parser.Consumer.Consume(LEFT_CURLY)!=nil {
		for parser.Consumer.Consume(RIGHT_CURLY)==nil {
			statements = append(statements, parser.Statement())
		}
		f.Body = statements
	} else {
		parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after fn deceleration")
	}
	f.Body = statements
	return f
}

// parse a variable definition (this only includes the identifier and type e.g. X : i32;, and assigning to
// a definition e.g. X : i32 = 1;)
func (parser *Parser) Define() AST {
	// explicit type define
	identifier := parser.Consumer.Consume(IDENTIFIER)
	parser.Consumer.Consume(COLON)
	def := &VarDefAST{
		Identifier: identifier,
		Type:       TavType{},
		Assignment: nil,
	}
	// check if we have a pointer
	def.Type = *parser.ParseType()
	// if the definition is of a struct or function, parse the body
	switch def.Type.Type {
	case TYPE_STRUCT:
		return parser.Struct(identifier)
	case TYPE_FN:
		return parser.Fn(identifier)
	default:
		if parser.Consumer.Consume(ASSIGN) != nil {
			def.Assignment = parser.Expression()
		}
		// add the identifier to the current symbol table
		parser.SymTable.Add(def.Identifier.Lexme(), def.Type, nil)
		return def
	}
}

// FOR NOW, WE DON'T SUPPORT QUICK ASSIGNING STRUCTS OR FUNCTIONS
// parse a variable quick assign (e.g. X := 1)
func (parser *Parser) QuickAssign() AST {
	identifier := parser.Consumer.Consume(IDENTIFIER)
	parser.Consumer.Consume(QUICK_ASSIGN)
	def := &VarDefAST{
		Identifier: identifier,
		Type:       TavType{},
		Assignment: nil,
	}
	// get the expression
	def.Assignment = parser.Expression()
	// then infer the type of the expression
	def.Type = InferType(def.Assignment, parser.SymTable)
	// add the identifier to the current symbol table
	parser.SymTable.Add(identifier.Lexme(), def.Type, nil)
	return def
}

// lowest precidence expression
func (parser *Parser) Expression() AST {
	return parser.Assignment()
}

func (parser *Parser) Assignment() AST {
	higherPrecedence := parser.ConnectiveOr()
	if parser.Consumer.Consume(ASSIGN) != nil {
		assignValue := parser.ConnectiveOr()
		// the only 2 types of assignments are to variables, and to struct members e.g. x = 2; or vec.x = 2;
		switch ast := higherPrecedence.(type) {
		case *VariableAST:
			return &VarSetAST{
				Identifier: ast.Identifier,
				Value:      assignValue,
			}
		case *StructGetAST:
			return &StructSetAST{
				Struct: ast.Struct,
				Member: ast.Member,
				Value:  assignValue,
				Deref:  ast.Deref,
			}
		}
	}
	return higherPrecedence
}

func (parser *Parser) ConnectiveOr() AST {
	higherPrecedence := parser.ConnectiveAnd()
	for parser.Consumer.Expect(OR) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.ConnectiveOr(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) ConnectiveAnd() AST {
	higherPrecedence := parser.BitwiseOr()
	for parser.Consumer.Expect(AND) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.ConnectiveAnd(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) BitwiseOr() AST {
	higherPrecedence := parser.BitwiseAnd()
	for parser.Consumer.Expect(BIN_OR) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.BitwiseOr(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) BitwiseAnd() AST {
	higherPrecedence := parser.Equality()
	for parser.Consumer.Expect(BIN_AND) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.BitwiseAnd(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) Equality() AST {
	higherPrecedence := parser.Comparison()
	for parser.Consumer.Expect(EQUALS) || parser.Consumer.Expect(NOT_EQUALS) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.Equality(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) Comparison() AST {
	higherPrecedence := parser.BitwiseShift()
	for parser.Consumer.Expect(GREAT_THAN) || parser.Consumer.Expect(GREAT_EQUAL) || parser.Consumer.Expect(LESS_THAN) || parser.Consumer.Expect(LESS_EQUAL) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.Comparison(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) BitwiseShift() AST {
	higherPrecedence := parser.PlusMinus()
	for parser.Consumer.Expect(SLEFT) || parser.Consumer.Expect(SRIGHT) {
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.BitwiseShift(),
		}
	}
	return higherPrecedence
}

func (parser *Parser) PlusMinus() AST {
	higherPrecedence := parser.MulDivModRem()
	for parser.Consumer.Expect(PLUS) || parser.Consumer.Expect(MINUS) {
		// TODO check here for compound assignment
		// if parser.Consumer.Expect(ASSIGN){...}
		b := &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.PlusMinus(),
		}
		return b
	}
	return higherPrecedence
}

func (parser *Parser) MulDivModRem() AST {
	higherPrecedence := parser.Unary()
	// TODO implement modulo and no-remainder division
	for parser.Consumer.Expect(STAR) || parser.Consumer.Expect(DIV) {
		// TODO check here for compound assignment
		// if parser.Consumer.Expect(ASSIGN){...}
		return &BinaryAST{
			Left:     higherPrecedence,
			Operator: parser.Consumer.Advance(),
			Right:    parser.MulDivModRem(),
		}
	}
	return higherPrecedence
}


func (parser *Parser) Unary() AST {
	// TODO implement increment and decrement
	// connective not, bitwise not, increment, decrement
	for parser.Consumer.Expect(BANG) || parser.Consumer.Expect(WIGGLE) {
		return &UnaryAST{
			Operator: parser.Consumer.Advance(),
			Right:    parser.Unary(),
		}
	}
	return parser.Casting()
}

func (parser *Parser) Casting() AST{
	if parser.Consumer.Expect(LEFT_PAREN) && parser.IsType(parser.Consumer.PeekAhead(1)) && parser.Consumer.ExpectAhead(RIGHT_PAREN, 2){
		Log("casting!")
		parser.Consumer.Consume(LEFT_PAREN)
		t := parser.ParseType()
		parser.Consumer.Consume(RIGHT_PAREN)
		return &CastAST{
			TavType: *t,
			Expr:    parser.Casting(),
		}
	}
	return parser.Addressing()
}

func (parser *Parser) Addressing() AST {
	if parser.Consumer.Expect(ADDR) || parser.Consumer.Expect(STAR) {
		return &UnaryAST{
			Operator: parser.Consumer.Advance(),
			Right:    parser.Addressing(),
		}
	}
	return parser.Call()
}

func (parser *Parser) Call() AST {
	callee := parser.SingleVal()
	// if the calle is a function e.g. 'main' and it doesn't have paramaters, it counts as a call
	// we need some way of check

	if InferType(callee, parser.SymTable).Type == TYPE_FN || InferType(callee, parser.SymTable).Type == TYPE_FN || parser.Consumer.Expect(LEFT_PAREN) {
		var args []AST
		if parser.Consumer.Consume(LEFT_PAREN) != nil {
			for !parser.Consumer.Expect(RIGHT_PAREN) {
				args = append(args, parser.Expression())
				if parser.Consumer.Expect(RIGHT_PAREN) {
					break
				}
				parser.Consumer.ConsumeErr(COMMA, ERR_UNEXPECTED_TOKEN, "expected ',' between arguments")
			}
			parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
		}
		return &CallAST{
			Caller: callee,
			Args:   args,
		}
	} else if parser.Consumer.Consume(PERIOD) != nil {
		// struct member get
		return &StructGetAST{
			Struct: callee,
			Member: parser.Consumer.ConsumeErr(IDENTIFIER, ERR_UNEXPECTED_TOKEN, "expected struct member"),
			Deref:  false,
		}
	} else if parser.Consumer.Consume(DEREF) != nil {
		// struct member dereference
		return &StructGetAST{
			Struct: callee,
			Member: parser.Consumer.ConsumeErr(IDENTIFIER, ERR_UNEXPECTED_TOKEN, "expected struct member"),
			Deref:  true,
		}
	}

	return callee
}

func (parser *Parser) SingleVal() AST {
	if t := parser.Consumer.Consume(IDENTIFIER); t != nil {
		return &VariableAST{Identifier: t}
	} else if t := parser.Consumer.Consume(NLITERAL); t != nil {
		// check if its a float
		if strings.Contains(t.Value.(string), ".") {
			value, err := strconv.ParseFloat(t.Value.(string), 64)
			Assert(err == nil, "couldn't parse float string value")
			return &LiteralAST{
				Type: TavType{
					Type: TYPE_F32,
				},
				Value: TavValue{
					Float: value,
				},
			}
		} else {
			value, err := strconv.ParseInt(t.Value.(string), 10, 64)
			Assert(err == nil, "couldn't parse int string value")
			return &LiteralAST{
				Type: TavType{
					Type: TYPE_I32,
				},
				Value: TavValue{
					Int: value,
				},
			}
		}
	} else if t := parser.Consumer.Consume(SLITERAL); t != nil {
		// parse the string literal into a character array here
		return &LiteralAST{
			Type: TavType{
				Type:        TYPE_STRING,
				// TODO this should probably be 1 as a string is just an *i8
				Indirection: 0,
				RetType:     nil,
			},
			Value: TavValue{
				String: []byte(t.Value.(string)+"\000"),
			},
		}
	} else if t := parser.Consumer.Consume(TRUE); t != nil {
		return &LiteralAST{
			Type: TavType{
				Type: TYPE_BOOL,
			},
			Value: TavValue{
				Bool: true,
			},
		}
	} else if t := parser.Consumer.Consume(FALSE); t != nil {
		return &LiteralAST{
			Type: TavType{
				Type: TYPE_BOOL,
			},
			Value: TavValue{
				Bool: false,
			},
		}
	} else if parser.Consumer.Consume(LEFT_PAREN) != nil { // group expression e.g. (1+2)
		expression := parser.Expression()
		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
		return &GroupAST{Group: expression}
	} else {
		parser.Compiler.Critical(parser.Consumer.Reporter, ERR_UNEXPECTED_TOKEN, "unexpected token")
	}
	return nil
}

// parse a type
func (parser *Parser) ParseType() *TavType {
	typ :=NewTavType(TYPE_VOID, "", 0, nil)
	// if it is a pointer, recursively get the pointer value
	for parser.Consumer.Consume(STAR) != nil {
		typ.Indirection += 1
	}
	if t := parser.Consumer.Consume(TYPE); t != nil {
		// it isn't a pointer, so get the type
		typ.Type = t.Value.(uint32)
	} else if t := parser.Consumer.Consume(IDENTIFIER); t != nil {
		typ.Type = TYPE_INSTANCE
		typ.Instance = t.Lexme()
	}
	return &typ
}


// returns true if the next token is a type
func (parser *Parser) IsType(token *Token) bool{
	return token.Type == IDENTIFIER || token.Type==TYPE
}