package src

import (
	"strconv"
	"strings"
)

const (
	ERR_UNEXPECTED_TOKEN = 0x0
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
	if parser.Consumer.Consume(RETURN) != nil {
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
		statements = append(statements, parser.Prelim())
	}
	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '}' at end of statement block")

	parser.SymTable = parser.SymTable.PopScope()

	return statements
}

// parse a struct
func (parser *Parser) Struct(identifier *Token) AST {
	s := &StructAST{Identifier: identifier}
	parser.Consumer.ConsumeErr(LEFT_CURLY, ERR_UNEXPECTED_TOKEN, "expected '{' after 'struct'")

	for !parser.Consumer.Expect(RIGHT_CURLY) {
		s.Fields = append(s.Fields, parser.Define().(*VarDefAST))
	}

	parser.Consumer.ConsumeErr(RIGHT_CURLY, ERR_UNEXPECTED_TOKEN, "expected closing '}'")
	return s
}

// parse a function
func (parser *Parser) Fun(identifier *Token) AST {
	f := &FnAST{Identifier: identifier, RetType: *parser.ParseType()}

	if parser.Consumer.Consume(LEFT_PAREN) != nil {

		// process the arguments

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
func (parser *Parser) Define() AST {
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
		// add the identifier to the current symbol table
		parser.SymTable.Add(identifier.Value.(string), def.Type, 0, nil)
		// if the definition is of a struct or function, parse the body
		switch def.Type.Type {
		case TYPE_STRUCT:
			return parser.Struct(identifier)
		case TYPE_FN:
			return parser.Fun(identifier)
		default:
			if parser.Consumer.Consume(ASSIGN) != nil {
				def.Assignment = parser.Expression()
			}
		}
	} else {
		parser.Consumer.Consume(QUICK_ASSIGN)
		// we are doing a quick assign
		assignment, assignmentType := parser.QuickAssign()
		def.Assignment = assignment
		def.Type = assignmentType
	}
	parser.Consumer.ConsumeErr(SEMICOLON, ERR_UNEXPECTED_TOKEN, "expected ';' after assignment")
	return def
}

// FOR NOW, WE DON'T SUPPORT QUICK ASSIGNING STRUCTS OR FUNCTIONS
// parse a variable quick assign (e.g. X := 1)
func (parser *Parser) QuickAssign() (AST, TavType) {
	// get the expression
	expression := parser.Expression()
	// then infer the type of the expression
	exprType := parser.InferType(expression)
	Log("quick assign type", exprType)
	return expression, exprType
}

func (parser *Parser) Expression() AST {

	return parser.Assignment()

	//parser.Consumer.Advance()
	//return &LiteralAST{
	//	Type:  TavType{
	//		Type:   TYPE_I32,
	//		IsPtr:  false,
	//		PtrVal: nil,
	//	},
	//	Value: TavValue{Int: 64},
	//}
}

// lowest precidence expression
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
	return parser.SingleVal()
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
	}
	return nil
}

// figure out the type of an expression
// TODO For now this will return ANY type
func (parser *Parser) InferType(expression AST) TavType {
	switch e := expression.(type){
	case *VariableAST:			// get the type of the variable in the symbol table
		return parser.SymTable.Get(e.Identifier.Value.(string)).Type
	case *LiteralAST:			// get the value of the literal
		return e.Type
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
