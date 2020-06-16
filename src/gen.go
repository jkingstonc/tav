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
	// the generator symbol table contains only value.Value
	SymTable     *SymTable
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
	generator.Module.NewTypeDef(StructAST.Identifier.Lexme(), s)
	return nil
}

// visit a function
// return the actual function
func (generator *Generator) VisitFnAST(FnAST *FnAST) interface{} {

	generator.SymTable = generator.SymTable.NewScope()

	f := generator.Module.NewFunc(FnAST.Identifier.Lexme(), LLType(FnAST.RetType))
	b := f.NewBlock("")
	generator.CurrentBlock = append(generator.CurrentBlock, b) // push the block to the stack
	for _, stmt := range FnAST.Body {
		stmt.Visit(generator)
	}
	generator.CurrentBlock = generator.CurrentBlock[:len(generator.CurrentBlock)-1] // pop the block from the stack

	generator.SymTable = generator.SymTable.PopScope()
	generator.SymTable.Add(FnAST.Identifier.Lexme(), TavType{
		Type:    TYPE_FN,
		IsPtr:   false,
		PtrVal:  nil,
		RetType: &FnAST.RetType,
	}, 0, f)
	return f
}

// visit a variable decleration
// return nothing
//	b.NewStore(assignment.(value.Value), v) as this is never used in a return evaulation
func (generator *Generator) VisitVarDefAST(VarDefAST *VarDefAST) interface{} {
	b := generator.CurrentBlock[len(generator.CurrentBlock)-1]
	// allocate memory on the stack & then store the assignment
	v := b.NewAlloca(LLType(VarDefAST.Type))
	if VarDefAST.Assignment != nil{
		assignment := VarDefAST.Assignment.Visit(generator)
		b.NewStore(assignment.(value.Value), v)
	}
	// store the allocated memory in the symbol table
	// whenever we want the variable value, we retrieve this
	generator.SymTable.Add(VarDefAST.Identifier.Lexme(), VarDefAST.Type, 0, v)
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
	// first get the variable from the symbol table (its a value.Value)
	variable := generator.SymTable.Get(VariableAST.Identifier.Lexme())
	b := generator.CurrentBlock[len(generator.CurrentBlock)-1]
	val := b.NewLoad(LLType(variable.Type), variable.Value.(value.Value))
	return val
}

func (generator *Generator) VisitUnaryAST(UnaryAST *UnaryAST) interface{} {
	switch UnaryAST.Operator.Type{
	case ADDR:
		b := generator.CurrentBlock[len(generator.CurrentBlock)-1]
		right := UnaryAST.Right.Visit(generator) // if this is a variable, it is a load instruction
		b.NewGetElementPtr(types.I32, right.(value.Value))
	}
	return nil
}

func (generator *Generator) VisitBinaryAST(BinaryAST *BinaryAST) interface{} {
	b := generator.CurrentBlock[len(generator.CurrentBlock)-1]
	left := BinaryAST.Left.Visit(generator).(value.Value)
	right := BinaryAST.Right.Visit(generator).(value.Value)
	switch BinaryAST.Operator.Type{
	case PLUS:
		if left.Type() == types.Float {
			return b.NewFAdd(left, right)
		}
		return b.NewAdd(left, right)
	case MINUS:
		if left.Type() == types.Float {
			return b.NewFSub(left, right)
		}
		return b.NewSub(left, right)
	case STAR:
		if left.Type() == types.Float {
			return b.NewFMul(left, right)
		}
		return b.NewMul(left, right)
	case DIV:
		return b.NewFDiv(left, right)
	}
	return nil
}

func (generator *Generator) VisitConnectiveAST(ConnectiveAST *ConnectiveAST) interface{} {
	return nil
}

func (generator *Generator) VisitCallAST(CallAST *CallAST) interface{} {
	b := generator.CurrentBlock[len(generator.CurrentBlock)-1]
	callee := CallAST.Caller.Visit(generator)
	return b.NewCall(callee.(value.Value))
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
	Log("generating")
	module := ir.NewModule()
	generator := &Generator{
		Root:   RootAST,
		Module: module,
		SymTable: NewSymTable(nil),
	}
	result := generator.Run()
	return result
}

func (generator *Generator) Run() *ir.Module {
	generator.Root.Visit(generator)
	return generator.Module
}
