package src

type Visitor interface {
	VisitStruct(s *StructAST) interface{}
	VisitFn(f *FnAST) interface{}
}

type AST interface {
}


// STATEMENTS

type RootAST struct {
	Statements    []AST
}

type ReturnAST struct{
	Value   AST
}
type ContinueAST struct{}
type BreakAST struct{}
type ForAST struct{

}
type IfAST struct{
	IfCondition   AST
	IfBody	      AST

	ElifCondition []AST
	ElifBody	  []AST

	ElseBody	  AST
}
type StructAST struct {
	Identifier *Token
}
type FnAST struct {
	Identifier  *Token
	Params      []VarDefAST	// the paramaters is an array of definitions
	Body     	*BlockAST
	RetType     uint32
	Variadic    bool
}
type VarDefAST struct{
	Identifier *Token
	Type 	   uint32
	Assignment AST
}
type BlockAST struct{
	Statements []AST
}
type ExprStmtAST struct{
	Expression AST
}


// EXPRESSIONS

type AsssignAST struct{			// assign to a non-member variable
	Identifier *Token
	Value	   AST
}
// Represents a literal value
type LiteralAST struct {
	Type uint32
	Value interface{}
}
type ListAST struct {}
type VariableAST struct {
	Identifier *Token
}
type UnaryAST struct {
	Operator  *Token
	Right     AST
}
type BinaryAST struct {
	Left	   AST
	Operator   *Token
	Right	   AST
}
type ConnectiveAST struct {
	Left	   AST
	Operator   *Token
	Right	   AST
}
type CallAST struct {
	Caller    AST
	Args      []AST
}
type StructGetAST struct {
	Struct    AST
	Member    *Token
}
type StructSetAST struct {
	Struct    AST
	Member    *Token
	Value     AST
}
type GroupAST struct {
	Group     AST
}