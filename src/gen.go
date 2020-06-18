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
	CurrentFn    *ir.Func
	// the generator symbol table contains only value.Value
	SymTable *SymTable
	Compiler *Compiler
}

func (Generator *Generator) PrintfProto() *ir.Func {
	f := Generator.Module.NewFunc("printf", types.I32, ir.NewParam("formatter", types.I8Ptr))
	Generator.SymTable.Add("printf", NewTavType(TYPE_FN, "", 0, nil), f)
	return f
}

func ValueFromType(tavType TavType, TavValue TavValue) value.Value {
	switch tavType.Type {
	case TYPE_BOOL:
		var val int64
		if TavValue.Bool == true {
			val = 1
		} else {
			val = 0
		}
		return constant.NewInt(types.I1, val)
	case TYPE_I8:
		return constant.NewInt(types.I8, TavValue.Int)
	case TYPE_I16:
		return constant.NewInt(types.I16, TavValue.Int)
	case TYPE_I32:
		return constant.NewInt(types.I32, TavValue.Int)
	case TYPE_I64:
		return constant.NewInt(types.I64, TavValue.Int)
	case TYPE_F32:
		return constant.NewFloat(types.Float, TavValue.Float)
	case TYPE_F64:
		return constant.NewFloat(types.Double, TavValue.Float)
	}
	return nil
}

func (generator *Generator) VisitRootAST(RootAST *RootAST) interface{} {
	generator.PrintfProto()
	for _, statement := range RootAST.Statements {
		statement.Visit(generator)
	}
	return nil
}

func (generator *Generator) VisitCastAST(CastAST *CastAST) interface{} {
	return nil
}

func (generator *Generator) VisitVarSetAST(VarSetAST *VarSetAST) interface{} {
	variable := generator.SymTable.Get(VarSetAST.Identifier.Lexme()).Value
	val := generator.Block().NewStore(VarSetAST.Value.Visit(generator).(value.Value), variable.(value.Value))
	return val
}

func (generator *Generator) VisitReturnAST(ReturnAST *ReturnAST) interface{} {
	generator.Block().NewRet(ReturnAST.Value.Visit(generator).(value.Value))
	return nil
}

func (generator *Generator) VisitBreakAST(BreakAST *BreakAST) interface{} {

	return nil
}

func (generator *Generator) VisitForAST(ForAST *ForAST) interface{} {
	return nil
}

func (generator *Generator) VisitIfAST(IfAST *IfAST) interface{} {

	ifCond:=generator.Block()
	ifBody:=generator.CurrentFn.NewBlock("if_body")

	var elifConditions []*ir.Block
	var elifBodies 	   []*ir.Block

	end:=generator.EnterBlock("end")



	// we need to see if we can cast the ifCOndition to a types.I1
	if len(elifConditions)>0 {
		ifCond.NewCondBr(IfAST.IfCondition.Visit(generator).(value.Value), ifBody, elifConditions[0])
	}else{
		ifCond.NewCondBr(IfAST.IfCondition.Visit(generator).(value.Value), ifBody, end)
	}


	ifBody.NewBr(end)
	for _, body := range elifBodies{
		body.NewBr(end)
	}

	// enter the if condition
	generator.Block().NewBr(ifCond)

	return nil
}

func (generator *Generator) VisitStructAST(StructAST *StructAST) interface{} {

	// create a new symbol table for the struct members
	generator.SymTable.NewScope(StructAST.Identifier.Lexme() + "_members")
	// create a new symbol table containing the children
	for _, member := range StructAST.Fields {
		generator.SymTable.Add(member.Identifier.Lexme(), member.Type, nil)
	}
	generator.SymTable.PopScope()

	// add the struct to the symbol table
	s := types.NewStruct()
	generator.SymTable.Add(StructAST.Identifier.Lexme(), NewTavType(TYPE_STRUCT, "", 0, nil), s)
	s.Packed = StructAST.Packed
	for _, field := range StructAST.Fields {
		s.Fields = append(s.Fields, ConvertType(field.Type, generator.SymTable))
	}
	generator.Module.NewTypeDef(StructAST.Identifier.Lexme(), s)
	return nil
}

func (generator *Generator) VisitFnAST(FnAST *FnAST) interface{} {

	identifier := FnAST.Identifier.Lexme()

	// add each function paramater to the function body scope
	generator.SymTable.NewScope(identifier + "_body")
	var params []*ir.Param
	for _, param := range FnAST.Params {
		p := ir.NewParam(param.Identifier.Lexme(), ConvertType(param.Type, generator.SymTable))
		params = append(params, p)
		generator.SymTable.Add(param.Identifier.Lexme(), param.Type, nil)
	}

	// create the function, the function body and visit the function block
	f := generator.Module.NewFunc(identifier, ConvertType(FnAST.RetType, generator.SymTable), params...)
	generator.CurrentFn = f
	b := f.NewBlock(identifier + "_body")
	generator.CurrentBlock = append(generator.CurrentBlock, b) // push the block to the stack
	for _, stmt := range FnAST.Body {
		stmt.Visit(generator)
	}
	generator.CurrentBlock = generator.CurrentBlock[:len(generator.CurrentBlock)-1] // pop the block from the stack

	generator.SymTable.PopScope()
	// finally add the function to the symbol table
	generator.SymTable.Add(identifier, NewTavType(TYPE_FN, "", 0, &FnAST.RetType), nil)
	return f
}

func (generator *Generator) VisitVarDefAST(VarDefAST *VarDefAST) interface{} {
	b := generator.Block()
	// allocate memory on the stack & then store the assignment
	v := b.NewAlloca(ConvertType(VarDefAST.Type, generator.SymTable))
	// if the variable assignment isn't nil, visit it and create an instruction to initialise the value
	if VarDefAST.Assignment != nil {
		assignment := VarDefAST.Assignment.Visit(generator)
		storeType := assignment.(value.Value)
		// if the type is a pointer or a struct, we have to use a load instruction before the store
		if VarDefAST.Type.Indirection != 0 || VarDefAST.Type.Type == TYPE_STRUCT {
			storeType = b.NewLoad(ConvertType(VarDefAST.Type, generator.SymTable), assignment.(value.Value))
		}
		b.NewStore(storeType, v)
	}
	// store the actual variable allocation in the symbol table
	// this will be retrieved any time we visit the variable
	generator.SymTable.Add(VarDefAST.Identifier.Lexme(), VarDefAST.Type, v)
	return nil
}

func (generator *Generator) VisitBlockAST(BlockAST *BlockAST) interface{} {
	generator.SymTable.NewScope("block_body")
	for _, stmt := range BlockAST.Statements {
		stmt.Visit(generator)
	}
	generator.SymTable.PopScope()
	return nil
}

func (generator *Generator) VisitExprSmtAST(ExprStmtAST *ExprStmtAST) interface{} {
	return ExprStmtAST.Expression.Visit(generator)
}

func (generator *Generator) VisitLiteralAST(LiteralAST *LiteralAST) interface{} {
	return ValueFromType(LiteralAST.Type, LiteralAST.Value)
}

func (generator *Generator) VisitListAST(ListAST *ListAST) interface{} {
	return nil
}

// return the value of a variable in the symbol table
// TODO This means functions are not first class variables as you cannot cast them to value.Value
func (generator *Generator) VisitVariableAST(VariableAST *VariableAST) interface{} {
	variable := generator.SymTable.Get(VariableAST.Identifier.Lexme())
	// when returning variables, we have to check the value
	// if the value is a function, we don't want to return a variable load instruction
	// instead we want to directly return the function to call
	// if the value is a paramater, we return the value directly
	switch variable.Type.Type {
	case TYPE_INSTANCE:
		return variable.Value
	case TYPE_FN:
		return variable.Value
	}
	switch variable.Value.(type) {
	case *ir.Param:
		return variable.Value
	default:
		val := generator.Block().NewLoad(ConvertType(variable.Type, generator.SymTable), variable.Value.(value.Value))
		return val
	}
}

func (generator *Generator) VisitUnaryAST(UnaryAST *UnaryAST) interface{} {
	b := generator.Block()
	right := UnaryAST.Right.Visit(generator) // if this is a variable, it is a load instruction
	// we have to invert the pointers (e.g. from *i32 -> i32 and vice versa)
	// the reason being the LLVM bindings want to see the target variable type
	// TODO multiple levels of indirection
	switch UnaryAST.Operator.Type {
	case ADDR:
		inverted := ConvertType(InvertPtrType(InferType(UnaryAST.Right, generator.SymTable), 1), generator.SymTable)
		//return	b.NewIntToPtr(right.(value.Value), inverted)
		return b.NewIntToPtr(right.(value.Value), inverted)
	case STAR:
		inverted := ConvertType(InvertPtrType(InferType(UnaryAST.Right, generator.SymTable), -1), generator.SymTable)
		//return b.NewPtrToInt(right.(value.Value), inverted)
		return b.NewPtrToInt(right.(value.Value), inverted)
	}
	return nil
}

func (generator *Generator) VisitBinaryAST(BinaryAST *BinaryAST) interface{} {
	b := generator.Block()
	left := BinaryAST.Left.Visit(generator).(value.Value)
	right := BinaryAST.Right.Visit(generator).(value.Value)
	switch BinaryAST.Operator.Type {
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
	callee := CallAST.Caller.Visit(generator)
	var args []value.Value
	for _, arg := range CallAST.Args {
		args = append(args, arg.Visit(generator).(value.Value))
	}
	return generator.Block().NewCall(callee.(value.Value), args...)
}

func (generator *Generator) VisitStructGetAST(StructGet *StructGetAST) interface{} {
	b := generator.Block()
	s := StructGet.Struct.Visit(generator)
	typeOfStruct := ConvertType(InferType(StructGet.Struct, generator.SymTable), generator.SymTable)
	member := b.NewGetElementPtr(typeOfStruct, s.(value.Value), constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	if !StructGet.Deref {
		// if we just return the member now, its a pointer! so we want to dereference it
		return b.NewPtrToInt(member, types.I32)
	}
	return member
}

func (generator *Generator) VisitStructSetAST(StructSetAST *StructSetAST) interface{} {
	b := generator.Block()
	s := StructSetAST.Struct.Visit(generator)
	typeOfStruct := ConvertType(InferType(StructSetAST.Struct, generator.SymTable), generator.SymTable)

	//offsett := CalcStructOffset(StructSetAST.Member.Lexme())

	member := b.NewGetElementPtr(typeOfStruct, s.(value.Value), constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	val := StructSetAST.Value.Visit(generator).(value.Value)
	b.NewStore(val, member)
	return member
}

// calculate the memory offset of a particular struct member
func (generator *Generator) CalcStructOffset(name, member string) int {
	// first find the struct in the symbol table
	t := generator.SymTable.Get(name).Value.(*Scope)
	for i, sym := range t.Symbols {
		if sym.Identifier == member {
			return i
		}
	}
	return 0
}

func (generator *Generator) VisitGroupAST(GroupAST *GroupAST) interface{} {
	return GroupAST.Group.Visit(generator)
}

func Generate(compiler *Compiler, RootAST *RootAST) *ir.Module {
	Log("generating")
	module := ir.NewModule()
	generator := &Generator{
		Root:     RootAST,
		Module:   module,
		SymTable: NewSymTable(),
		Compiler: compiler,
	}
	result := generator.Run()
	return result
}

func (generator *Generator) Run() *ir.Module {
	generator.Root.Visit(generator)
	return generator.Module
}

// generate a new generator scope by pushing a new symbol table scope and a new block
func (generator *Generator) NewBlock(identifier string) *ir.Block{
	block := generator.CurrentFn.NewBlock(identifier)
	return block
}

func (generator *Generator) EnterBlock(identifier string) *ir.Block{
	generator.CurrentBlock = append(generator.CurrentBlock, generator.NewBlock(identifier))
	return generator.Block()
}

func (generator *Generator) ExitBlock(){
	generator.CurrentBlock = generator.CurrentBlock[:len(generator.CurrentBlock)-1] // pop the block from the stack
}

func (generator *Generator) Block() *ir.Block{
	return generator.CurrentBlock[len(generator.CurrentBlock)-1]
}