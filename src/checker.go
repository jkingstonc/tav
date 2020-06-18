package src

const (
	ERR_REDECLARED          = 0x0
	ERR_INVALID_RETURN_TYPE = 0x1
	ERR_NO_VAR              = 0x2
)

// implements Visitor
type Checker struct {
	Compiler *Compiler
	// used for scope checking
	SymTable *SymTable
	Reporter *Reporter
	Root     *RootAST
}

func (Checker *Checker) PrintfProto() {
	Checker.SymTable.Add("printf", TavType{
		Type:        TYPE_I8,
		Indirection: 1,
		RetType:     &TavType{
			Type:        TYPE_I32,
			Indirection: 0,
			RetType:     nil,
		},
	}, 0, nil, nil)
}

func Check(compiler *Compiler, RootAST *RootAST) *RootAST {
	reporter := NewReporter(compiler.File.Filename, compiler.File.Source)
	checker := Checker{
		Compiler: compiler,
		SymTable: NewSymTable(nil),
		Reporter: reporter,
		Root:     RootAST,
	}
	checker.Run()
	return RootAST
}

func (checker *Checker) Run() {
	checker.Root.Visit(checker)
}

func (checker *Checker) VisitRootAST(RootAST *RootAST) interface{} {
	checker.PrintfProto()
	for _, statement := range RootAST.Statements {
		statement.Visit(checker)
	}
	return nil
}

func (checker *Checker) VisitCastAST(CastAST *CastAST) interface{}{
	return nil
}

func (checker *Checker) VisitVarSetAST(VarSetAST *VarSetAST) interface{} {
	checker.Reporter.Position = VarSetAST.Identifier.Position
	if InferType(VarSetAST.Value, checker.SymTable) != checker.SymTable.Get(VarSetAST.Identifier.Lexme()).Type {
		checker.Compiler.Critical(checker.Reporter, ERR_INVALID_TYPE, "cannot assign type to variable")
	}
	return nil
}

func (checker *Checker) VisitReturnAST(ReturnAST *ReturnAST) interface{} {
	ReturnAST.Value.Visit(checker)
	return nil
}

func (checker *Checker) VisitBreakAST(BreakAST *BreakAST) interface{} {
	return nil
}

func (checker *Checker) VisitForAST(ForAST *ForAST) interface{} {
	return nil
}

func (checker *Checker) VisitIfAST(IfAST *IfAST) interface{} {
	return nil
}

func (checker *Checker) VisitStructAST(StructAST *StructAST) interface{} {

	checker.SymTable = checker.SymTable.NewScope()
	memberSymTable := checker.SymTable
	// create a new symbol table containing the children
	for _, member := range StructAST.Fields{
		checker.SymTable.Add(member.Identifier.Lexme(), member.Type, 0, nil, nil)
	}
	checker.SymTable = checker.SymTable.PopScope()

	checker.Reporter.Position = StructAST.Identifier.Position
	checker.SymTable.Add(StructAST.Identifier.Lexme(), TavType{
		Type:        TYPE_STRUCT,
	}, 0, nil, memberSymTable)
	return nil
}

func (checker *Checker) VisitFnAST(FnAST *FnAST) interface{} {
	checker.Reporter.Position = FnAST.Identifier.Position

	if checker.SymTable.Get(FnAST.Identifier.Lexme()) != nil {
		checker.Compiler.Critical(checker.Reporter, ERR_REDECLARED, "function re-declared")
	}
	checker.SymTable.Add(FnAST.Identifier.Lexme(), TavType{
		Type:    TYPE_FN,
		RetType: &FnAST.RetType,
	}, 0, nil, nil)

	checker.SymTable = checker.SymTable.NewScope()

	for _, param := range FnAST.Params {
		param.Visit(checker)
	}

	for _, stmt := range FnAST.Body {
		stmt.Visit(checker)
		switch s := stmt.(type) {
		case *ReturnAST:
			// check if the return value is of the same type
			if InferType(s.Value, checker.SymTable) != FnAST.RetType {
				checker.Reporter.Position = FnAST.Identifier.Position
				// cast the value to the return value automatically
				if !CastValue(FnAST.RetType, s.Value){
					checker.Compiler.Critical(checker.Reporter, ERR_INVALID_RETURN_TYPE, "return types do not match")
				}
			}
		}
	}
	checker.SymTable = checker.SymTable.PopScope()
	return nil
}

// visit a variable decleration
// return nothing
//	b.NewStore(assignment.(value.Value), v) as this is never used in a return evaulation
func (checker *Checker) VisitVarDefAST(VarDefAST *VarDefAST) interface{} {
	checker.Reporter.Position = VarDefAST.Identifier.Position

	// only check the local scope, otherwise we can redeclare global variables
	if checker.SymTable.GetLocal(VarDefAST.Identifier.Lexme()) != nil {
		checker.Compiler.Critical(checker.Reporter, ERR_REDECLARED, "variable re-declared")
	}
	// add the define to the symbol table
	checker.SymTable.Add(VarDefAST.Identifier.Lexme(), VarDefAST.Type, 0, nil, nil)
	// check if the assigned type was correct
	if VarDefAST.Assignment != nil {
		if t := InferType(VarDefAST.Assignment, checker.SymTable); t != VarDefAST.Type {
			if !CastValue(VarDefAST.Type, VarDefAST.Assignment){
				checker.Compiler.Critical(checker.Reporter, ERR_INVALID_TYPE, "types do not match")
			}
		}
	}
	return nil
}

func (checker *Checker) VisitBlockAST(BlockAST *BlockAST) interface{} {
	checker.SymTable = checker.SymTable.NewScope()
	for _, stmt := range BlockAST.Statements {
		stmt.Visit(checker)
	}
	checker.SymTable = checker.SymTable.PopScope()
	return nil
}

func (checker *Checker) VisitExprSmtAST(ExprStmtAST *ExprStmtAST) interface{} {
	ExprStmtAST.Visit(checker)
	return nil
}

func (checker *Checker) VisitLiteralAST(LiteralAST *LiteralAST) interface{} {
	return nil
}

func (checker *Checker) VisitListAST(ListAST *ListAST) interface{} {
	return nil
}

func (checker *Checker) VisitVariableAST(VariableAST *VariableAST) interface{} {
	checker.Reporter.Position = VariableAST.Identifier.Position
	if checker.SymTable.Get(VariableAST.Identifier.Lexme()) == nil {
		checker.Compiler.Critical(checker.Reporter, ERR_NO_VAR, "variable doesn't exist")
	}
	return nil
}

func (checker *Checker) VisitUnaryAST(UnaryAST *UnaryAST) interface{} {
	return nil
}

func (checker *Checker) VisitBinaryAST(BinaryAST *BinaryAST) interface{} {
	return nil
}

func (checker *Checker) VisitConnectiveAST(ConnectiveAST *ConnectiveAST) interface{} {
	return nil
}

func (checker *Checker) VisitCallAST(CallAST *CallAST) interface{} {
	CallAST.Caller.Visit(checker)
	for _, arg := range CallAST.Args {
		arg.Visit(checker)
	}
	return nil
}

func (checker *Checker) VisitStructGetAST(StructGet *StructGetAST) interface{} {
	return InferType(StructGet, checker.SymTable)
}

func (checker *Checker) VisitStructSetAST(StructSetAST *StructSetAST) interface{} {
	return nil
}

func (checker *Checker) VisitGroupAST(GroupAST *GroupAST) interface{} {
	return nil
}
