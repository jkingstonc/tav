package src

import (
	"github.com/llir/llvm/ir/types"
	"time"

	"github.com/llir/llvm/ir"
)

// implements visitor
type Generator struct {
	Root   *RootAST
	Module *ir.Module
}

func (generator *Generator) VisitRootAST(RootAST *RootAST) {
	for _, statement := range RootAST.Statements{
		statement.Visit(generator)
	}
}

func (generator *Generator) VisitReturnAST(ReturnAST *ReturnAST) {
}

func (generator *Generator) VisitBreakAST(BreakAST *BreakAST) {
}

func (generator *Generator) VisitForAST(ForAST *ForAST) {
}

func (generator *Generator) VisitIfAST(IfAST *IfAST) {
}

func (generator *Generator) VisitStructAST(StructAST *StructAST) {
}

func (generator *Generator) VisitFnAST(FnAST *FnAST) {
	generator.Module.NewFunc(FnAST.Identifier.Value.(string), types.I32)
}

func (generator *Generator) VisitVarDefAST(VarDefAST *VarDefAST) {
}

func (generator *Generator) VisitBlockAST(BlockAST *BlockAST) {
}

func (generator *Generator) VisitExprSmtAST(ExprStmtAST *ExprStmtAST) {
}

func (generator *Generator) VisitAssignAST(AsssignAST *AsssignAST) {
}

func (generator *Generator) VisitLiteralAST(LiteralAST *LiteralAST) {
}

func (generator *Generator) VisitListAST(ListAST *ListAST) {
}

func (generator *Generator) VisitVariableAST(VariableAST *VariableAST) {
}

func (generator *Generator) VisitUnaryAST(UnaryAST *UnaryAST) {
}

func (generator *Generator) VisitBinaryAST(BinaryAST *BinaryAST) {
}

func (generator *Generator) VisitConnectiveAST(ConnectiveAST *ConnectiveAST) {
}

func (generator *Generator) VisitCallAST(CallAST *CallAST) {
}

func (generator *Generator) VisitStructGetAST(StructGet *StructGetAST) {
}

func (generator *Generator) VisitStructSetAST(StructSetAST *StructSetAST) {
}

func (generator *Generator) VisitGroupAST(GroupAST *GroupAST) {
}

func Generate(compiler *Compiler, RootAST *RootAST) uint8 {
	start := time.Now()

	// create a new LLVM module
	module := ir.NewModule()
	generator := &Generator{
		Root:   RootAST,
		Module: module,
	}

	result := generator.Run()

	Log("generator result")
	Log(generator.Module.String())

	end := time.Since(start)
	Log("back end took ", end.Seconds(), "seconds")
	return result
}

func (generator *Generator) Run() uint8 {

	generator.Root.Visit(generator)

	return SUCCESS_COMP
}