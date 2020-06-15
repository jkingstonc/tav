package src

import (
	"strconv"
	"strings"
)

const (
	ERR_UNEXPECTED_TOKEN   = 0x0
	ERR_INVALID_TYPE 	   = 0x1
	ERR_INVALID_IDENTIFIER = 0x2
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
	reporter := NewReporter(compiler.FileName, compiler.Source)
	consumer := NewParseConsumer(tokens, reporter, compiler)

	parser := Parser{
		Compiler: compiler,
		Consumer: consumer,
		SymTable: NewSymTable(nil),
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
			Root.Statements = append(Root.Statements, parser.Define(true))
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
	if parser.Consumer.Expect(IDENTIFIER){
		return parser.Define(true)
	} else if parser.Consumer.Consume(RETURN) != nil {
		return parser.Return()
	} else if parser.Consumer.Consume(BREAK) != nil {
		return parser.Break()
	} else if parser.Consumer.Consume(FOR) != nil {
		return parser.For()
	} else if parser.Consumer.Consume(IF) != nil {
		return parser.If()
	} else if parser.Consumer.Consume(LEFT_CURLY) != nil {
		return &BlockAST{Statements: parser.ParseStmtBlock()}
	} else {
		return parser.ExpressionStmt()
	}
}

func (parser *Parser) ExpressionStmt() AST {
	return &ExprStmtAST{Expression: parser.Expression()}
}

func (parser *Parser) Return() AST {
	r := &ReturnAST{Value: parser.Expression()}
	parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after 'return'")
	return r
}

func (parser *Parser) Break() AST {
	r := &ReturnAST{Value: parser.Expression()}
	return r
}

func (parser *Parser) For() AST {
	r := &ReturnAST{Value: parser.Expression()}
	return r
}

func (parser *Parser) If() AST {
	r := &ReturnAST{Value: parser.Expression()}
	return r
}

func (parser *Parser) ParseStmtBlock() []AST {

	parser.SymTable = parser.SymTable.NewScope()

	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' at start of statement block")
	var statements []AST
	for !parser.Consumer.Expect(RIGHT_CURLY) && !parser.Consumer.End() {
		statements = append(statements, parser.Statement())
	}
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}' at end of statement block")

	parser.SymTable = parser.SymTable.PopScope()

	return statements
}

// parse a struct
func (parser *Parser) Struct(identifier *Token) AST {
	// add the identifier to the current symbol table
	parser.SymTable.Add(identifier.Value.(string), TavType{
		Type:   TYPE_STRUCT,
		IsPtr:  false,
		PtrVal: nil,
	}, 0, nil)
	s := &StructAST{Identifier: identifier}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' after 'struct'")

	for !parser.Consumer.Expect(RIGHT_CURLY) {
		s.Fields = append(s.Fields, parser.Define(true).(*VarDefAST))
	}

	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	return s
}

// parse a function
func (parser *Parser) Fn(identifier *Token) AST {	// add the identifier to the current symbol table
	f := &FnAST{Identifier: identifier, RetType: *parser.ParseType()}

	// we should probably instead use the return type???
	parser.SymTable.Add(identifier.Value.(string), TavType{
		Type:   TYPE_FN,
		IsPtr:  false,
		PtrVal: nil,
		RetType: &f.RetType,
	},  0, f)

	if parser.Consumer.Consume(LEFT_PAREN) != nil {
		var params []*VarDefAST
		// process the arguments
		for !parser.Consumer.Expect(RIGHT_PAREN) {
			// each paramater is essentially a variable decleration
			params = append(params, parser.Define(false).(*VarDefAST))
			if parser.Consumer.Expect(RIGHT_PAREN){
				break
			}
			parser.Consumer.ConsumeErr(COMMA, ERR_UNEXPECTED_TOKEN, "expected ',' between paramaters")
		}
		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
	}
	if parser.Consumer.Expect(LEFT_CURLY) {
		statements := parser.ParseStmtBlock()
		f.Body = statements
	} else {
		parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after fn decleration")
	}
	return f
}

// parse a variable definition (this only includes the identifier and type e.g. X : i32;, and assigning to
// a definition e.g. X : i32 = 1;)
func (parser *Parser) Define(expectSemiColon bool) AST {
	identifier := parser.Consumer.Consume(IDENTIFIER)

	def := &VarDefAST{
		Identifier: identifier,
		Type:       TavType{},
		Assignment: nil,
	}
	// explicit type define
	if parser.Consumer.Consume(COLON) != nil {
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
				// check if the assigned type was correct
				if parser.InferType(def.Assignment) != def.Type {
					parser.Compiler.Critical(parser.Consumer.Reporter, ERR_INVALID_TYPE, "invalid type")
				}
			}
		}
	} else {
		parser.Consumer.Consume(QUICK_ASSIGN)
		// we are doing a quick assign
		assignment, assignmentType := parser.QuickAssign()
		def.Assignment = assignment
		def.Type = assignmentType
	}
	// add the identifier to the current symbol table
	parser.SymTable.Add(identifier.Value.(string), def.Type, 0, nil)
	if expectSemiColon {
		parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after definition")
	}
	return def
}

// FOR NOW, WE DON'T SUPPORT QUICK ASSIGNING STRUCTS OR FUNCTIONS
// parse a variable quick assign (e.g. X := 1)
func (parser *Parser) QuickAssign() (AST, TavType) {
	// get the expression
	expression := parser.Expression()
	// then infer the type of the expression
	exprType := parser.InferType(expression)
	return expression, exprType
}

// lowest precidence expression
func (parser *Parser) Expression() AST {
	return parser.Assignment()
}

func (parser *Parser) Assignment() AST{
	higherPrecedence := parser.ConnectiveOr()
	return higherPrecedence
}

func (parser *Parser) ConnectiveOr() AST{
	return parser.ConnectiveAnd()
}

func (parser *Parser) ConnectiveAnd() AST{
	return parser.BitwiseOr()
}

func (parser *Parser) BitwiseOr() AST{
	return parser.BitwiseAnd()
}

func (parser *Parser) BitwiseAnd() AST{
	return parser.Equality()
}

func (parser *Parser) Equality() AST{
	return parser.Comparison()
}

func (parser *Parser) Comparison() AST{
	return parser.BitwiseShift()
}

func (parser *Parser) BitwiseShift() AST{
	return parser.PlusMinus()
}

func (parser *Parser) PlusMinus() AST{
	return parser.MulDivModRem()
}

func (parser *Parser) MulDivModRem() AST{
	return parser.Unary()
}

func (parser *Parser) Unary() AST{
	return parser.Call()
}

func (parser *Parser) Call() AST{
	callee := parser.SingleVal()
	// if the calle is a function e.g. 'main' and it doesn't have paramaters, it counts as a call
	// we need some way of check
	if parser.InferType(callee).Type == TYPE_FN {
		var args []AST
		if parser.Consumer.Consume(LEFT_PAREN) != nil{
			for !parser.Consumer.Expect(RIGHT_PAREN) {
				args = append(args, parser.Expression())
				if parser.Consumer.Expect(RIGHT_PAREN){
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
	}

	return callee
}

func (parser *Parser) SingleVal() AST{
	if t:=parser.Consumer.Consume(IDENTIFIER); t!=nil{
		return &VariableAST{Identifier: t}
	}else if t:=parser.Consumer.Consume(NLITERAL); t!=nil{
		// check if its a float
		if strings.Contains(t.Value.(string), ".") {
			value, err := strconv.ParseFloat(t.Value.(string), 64)
			Assert(err==nil, "couldn't parse float string value")
			return &LiteralAST{
				Type: TavType{
					Type:   TYPE_F32,
					IsPtr:  false,
					PtrVal: nil,
				},
				Value: TavValue{
					Int:    0,
					Float:  value,
					String: "",
					Bool:   false,
					Any:    nil,
				},
			}
		}else{
			value, err := strconv.ParseInt(t.Value.(string), 10, 64)
			Assert(err==nil, "couldn't parse int string value")
			return &LiteralAST{
				Type: TavType{
					Type:   TYPE_I32,
					IsPtr:  false,
					PtrVal: nil,
				},
				Value: TavValue{
					Int:    value,
					Float:  0,
					String: "",
					Bool:   false,
					Any:    nil,
				},
			}
		}
	}else if t:=parser.Consumer.Consume(SLITERAL);t!=nil{
		return &LiteralAST{
			Type: TavType{
				Type:   TYPE_STRING,
				IsPtr:  false,
				PtrVal: nil,
			},
			Value: TavValue{
				Int:    0,
				Float:  0,
				String: t.Value.(string),
				Bool:   false,
				Any:    nil,
			},
		}
	}else if t:=parser.Consumer.Consume(BOOL);t!=nil{
		return &LiteralAST{
			Type: TavType{
				Type:   TYPE_BOOL,
				IsPtr:  false,
				PtrVal: nil,
			},
			Value: TavValue{
				Int:    0,
				Float:  0,
				String: "",
				Bool:   t.Value.(bool),
				Any:    nil,
			},
		}
	}else if parser.Consumer.Consume(LEFT_PAREN) != nil{ // group expression e.g. (1+2)
		expression := parser.Expression()
		parser.Consumer.ConsumeErr(RIGHT_PAREN, ERR_UNEXPECTED_TOKEN, "expected closing ')'")
		return &GroupAST{Group: expression}
	}
	return nil
}

// figure out the type of an expression
// TODO For now this will return ANY type
func (parser *Parser) InferType(expression AST) TavType {
	switch e := expression.(type){
	case *VariableAST:			// get the type of the variable in the symbol table
		t := parser.SymTable.Get(e.Identifier.Value.(string))
		if t != nil{
			return t.Type
		}
		parser.Compiler.Critical(parser.Consumer.Reporter, ERR_INVALID_IDENTIFIER, "identifier doesn't exist")
	case *LiteralAST:			// get the value of the literal
		return e.Type
	case *CallAST:
		// TODO figure out how we infer the type of a function call
		t := parser.InferType(e.Caller)
		Log("infering type of CallAST...", t.RetType)
		return *t.RetType
	}
	// this is unreachable
	return TavType{}
}

// join 2 infered types and figure out what the next type will be
func (parser *Parser) JoinInfered(type1, type2 TavType) TavType {
	return type1
}

// parse a type
func (parser *Parser) ParseType() *TavType {
	typ := &TavType{
		Type:   0,
		IsPtr:  false,
		PtrVal: nil,
	}
	// if it is a pointer, recursively get the pointer value
	if parser.Consumer.Consume(STAR) != nil {
		typ.IsPtr = true
		typ.PtrVal = parser.ParseType()
	} else {
		if t := parser.Consumer.Consume(TYPE); t != nil {
			// it isn't a pointer, so get the type
			typ.Type = t.Value.(uint32)
		}
	}
	return typ
}
