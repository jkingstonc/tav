package src

type ASTVisitor interface {
	VisitStruct() *AST
	VisitFun() *AST
	VisitExpr() *AST
}

type AST struct {

}

type StructAST struct {
	AST
	Identifier *Token
}

type FunAST struct {
	AST
	Identifier *Token
}

type ExprAST struct {

}