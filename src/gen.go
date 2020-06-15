package src

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// implements visitor
type Generator struct {
	Root         *RootAST
	Module       *ir.Module
	CurrentBlock []*ir.Block
}

func ValueFromType(tavType TavType, TavValue TavValue) value.Value {
	llType := LLType(tavType)
	switch llType {
	case types.I1:
		var val int64
		if TavValue.Bool == true{
			val = 1
		}else{
			val = 0
		}
		return constant.NewInt(types.I1, val)
	case types.I8:
		return constant.NewInt(types.I8, TavValue.Int)
	case types.I16:
		return constant.NewInt(types.I16, TavValue.Int)
	case types.I32:
		return constant.NewInt(types.I32, TavValue.Int)
	case types.I64:
		return constant.NewInt(types.I64, TavValue.Int)
	case types.Float:
		return constant.NewFloat(types.Float, TavValue.Float)
	case types.Double:
		return constant.NewFloat(types.Double, TavValue.Float)
	}
	return nil
}

func (generator *Generator) VisitRootAST(RootAST *RootAST) interface{} {
	for _, statement := range RootAST.Statements {
		statement.Visit(generator)
	}
	return nil
}

func (generator *Generator) VisitVarSetAST(VarSetAST *VarSetAST) interface{} {
	return nil
}

func (generator *Generator) VisitReturnAST(ReturnAST *ReturnAST) interface{} {
	b := generator.CurrentBlock[len(generator.CurrentBlock)-1]
	b.NewRet(ReturnAST.Value.Visit(generator).(value.Value))
	return nil
}

func (generator *Generator) VisitBreakAST(BreakAST *BreakAST) interface{} {
	return nil
}

func (generator *Generator) VisitForAST(ForAST *ForAST) interface{} {
	return nil
}

func (generator *Generator) VisitIfAST(IfAST *IfAST) interface{} {
	return nil
}

func (generator *Generator) VisitStructAST(StructAST *StructAST) interface{} {
	s := types.NewStruct()
	s.Packed = StructAST.Packed
	for _, field := range StructAST.Fields {
		s.Fields = append(s.Fields, LLType(field.Type))
	}
	generator.Module.NewTypeDef(StructAST.Identifier.Value.(string), s)
	return nil
}

func (generator *Generator) VisitFnAST(FnAST *FnAST) interface{} {
	f := generator.Module.NewFunc(FnAST.Identifier.Value.(string), LLType(FnAST.RetType))
	b := f.NewBlock("")
	generator.CurrentBlock = append(generator.CurrentBlock, b) // push the block to the stack
	for _, stmt := range FnAST.Body {
		stmt.Visit(generator)
	}
	generator.CurrentBlock = generator.CurrentBlock[:len(generator.CurrentBlock)-1] // pop the block from the stack
	return nil
}

func (generator *Generator) VisitVarDefAST(VarDefAST *VarDefAST) interface{} {
	return nil
}

func (generator *Generator) VisitBlockAST(BlockAST *BlockAST) interface{} {
	return nil
}

func (generator *Generator) VisitExprSmtAST(ExprStmtAST *ExprStmtAST) interface{} {
	return nil
}

func (generator *Generator) VisitAssignAST(AsssignAST *AsssignAST) interface{} {
	return nil
}

func (generator *Generator) VisitLiteralAST(LiteralAST *LiteralAST) interface{} {
	return ValueFromType(LiteralAST.Type, LiteralAST.Value)
}

func (generator *Generator) VisitListAST(ListAST *ListAST) interface{} {
	return nil
}

func (generator *Generator) VisitVariableAST(VariableAST *VariableAST) interface{} {
	return nil
}

func (generator *Generator) VisitUnaryAST(UnaryAST *UnaryAST) interface{} {
	return nil
}

func (generator *Generator) VisitBinaryAST(BinaryAST *BinaryAST) interface{} {
	return nil
}

func (generator *Generator) VisitConnectiveAST(ConnectiveAST *ConnectiveAST) interface{} {
	return nil
}

func (generator *Generator) VisitCallAST(CallAST *CallAST) interface{} {
	return nil
}

func (generator *Generator) VisitStructGetAST(StructGet *StructGetAST) interface{} {
	return nil
}

func (generator *Generator) VisitStructSetAST(StructSetAST *StructSetAST) interface{} {
	return nil
}

func (generator *Generator) VisitGroupAST(GroupAST *GroupAST) interface{} {
	return nil
}

func Generate(compiler *Compiler, RootAST *RootAST) *ir.Module {
	module := ir.NewModule()
	generator := &Generator{
		Root:   RootAST,
		Module: module,
	}
	result := generator.Run()
	return result
}

func (generator *Generator) Run() *ir.Module {
	generator.Root.Visit(generator)
	return generator.Module
}
