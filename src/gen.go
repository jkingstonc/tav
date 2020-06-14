package src

import (
	"github.com/llir/llvm/ir/constant"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// implements visitor
type Generator struct {
	Root   		 *RootAST
	Module 	     *ir.Module
	CurrentBlock []*ir.Block
}

func (generator *Generator) VisitRootAST(RootAST *RootAST) {
	for _, statement := range RootAST.Statements{
		statement.Visit(generator)
	}
}

func (generator *Generator) VisitReturnAST(ReturnAST *ReturnAST) {
	b:=generator.CurrentBlock[len(generator.CurrentBlock)-1]
	b.NewRet(constant.NewInt(types.I32, 1))
}

func (generator *Generator) VisitBreakAST(BreakAST *BreakAST) {
}

func (generator *Generator) VisitForAST(ForAST *ForAST) {
}

func (generator *Generator) VisitIfAST(IfAST *IfAST) {
}

func (generator *Generator) VisitStructAST(StructAST *StructAST) {
	s := types.NewStruct()
	s.Packed = StructAST.Packed
	for _, field := range StructAST.Fields{
		s.Fields = append(s.Fields, LLType(field.Type))
	}
	generator.Module.NewTypeDef(StructAST.Identifier.Value.(string), s)
}

func (generator *Generator) VisitFnAST(FnAST *FnAST) {
	f := generator.Module.NewFunc(FnAST.Identifier.Value.(string), LLType(FnAST.RetType))
	b := f.NewBlock("")
	generator.CurrentBlock = append(generator.CurrentBlock, b)		// push the block to the stack
	for _, stmt := range FnAST.Body{
		stmt.Visit(generator)
	}
	generator.CurrentBlock = generator.CurrentBlock[:len(generator.CurrentBlock) - 1] // pop the block from the stack
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
	Log(generator.Module.String())

	ioutil.WriteFile("tmp/test.ll", []byte(generator.Module.String()), 0644)
	c:=exec.Command("llc","tmp/test.ll")
	err := c.Run()
	Log(err)
	c=exec.Command("gcc","-c","tmp/test.s","-o", "tmp/test.o")
	err = c.Run()
	Log(err)
	c=exec.Command("gcc", "tmp/test.o", "-o", "tmp/test")
	err = c.Run()
	Log(err)

	end := time.Since(start)
	Log("back end took ", end.Seconds(), "seconds")
	return result
}

func (generator *Generator) Run() uint8 {

	generator.Root.Visit(generator)

	return SUCCESS_COMP
}