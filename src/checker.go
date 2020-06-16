package src

const (
	ERR_REDECLARED = 0x0
	ERR_INVALID_RETURN_TYPE = 0x1
)

// implements Visitor
type Checker struct {
	Compiler *Compiler
	// used for scope checking
	SymTable  *SymTable
	Reporter *Reporter
	Root *RootAST
}

func Check (compiler *Compiler, RootAST *RootAST) *RootAST{
	reporter := NewReporter(compiler.FileName, compiler.Source)
	checker := Checker{
		Compiler: compiler,
		SymTable: NewSymTable(nil),
		Reporter: reporter,
		Root: RootAST,
	}
	checker.Run()
	return RootAST
}

func (checker *Checker) Run() {
	checker.Root.Visit(checker)
}

func (checker *Checker) VisitRootAST(RootAST *RootAST) interface{} {
	for _, statement := range RootAST.Statements {
		statement.Visit(checker)
	}
	return nil
}

func (checker *Checker) VisitVarSetAST(VarSetAST *VarSetAST) interface{} {
	return nil
}

func (checker *Checker) VisitReturnAST(ReturnAST *ReturnAST) interface{} {
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
	return nil
}

func (checker *Checker) VisitFnAST(FnAST *FnAST) interface{} {
	checker.Reporter.Position = FnAST.Identifier.Position
	checker.SymTable = checker.SymTable.NewScope()
	for _, stmt := range FnAST.Body{
		stmt.Visit(checker)
		switch s:=stmt.(type){
		case *ReturnAST:
			// check if the return value is of the same type
			if InferType(s.Value, checker.SymTable) != FnAST.RetType{
				checker.Reporter.Position = FnAST.Identifier.Position
				checker.Compiler.Critical(checker.Reporter, ERR_INVALID_RETURN_TYPE, "return types do not match")
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

	if checker.SymTable.Get(VarDefAST.Identifier.Lexme())!=nil{
		checker.Compiler.Critical(checker.Reporter, ERR_REDECLARED, "variable re-declared")
	}
	// add the define to the symbol table
	checker.SymTable.Add(VarDefAST.Identifier.Lexme(), VarDefAST.Type, 0, nil)
	// check if the assigned type was correct
	if t :=InferType(VarDefAST.Assignment, checker.SymTable); t != VarDefAST.Type {
		checker.Compiler.Critical(checker.Reporter, ERR_INVALID_TYPE, "types do not match")
	}
	return nil
}

func (checker *Checker) VisitBlockAST(BlockAST *BlockAST) interface{} {
	return nil
}

func (checker *Checker) VisitExprSmtAST(ExprStmtAST *ExprStmtAST) interface{} {
	return nil
}

func (checker *Checker) VisitAssignAST(AsssignAST *AsssignAST) interface{} {
	return nil
}

func (checker *Checker) VisitLiteralAST(LiteralAST *LiteralAST) interface{} {
	return nil
}

func (checker *Checker) VisitListAST(ListAST *ListAST) interface{} {
	return nil
}

func (checker *Checker) VisitVariableAST(VariableAST *VariableAST) interface{} {
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
	return nil
}

func (checker *Checker) VisitStructGetAST(StructGet *StructGetAST) interface{} {
	return nil
}

func (checker *Checker) VisitStructSetAST(StructSetAST *StructSetAST) interface{} {
	return nil
}

func (checker *Checker) VisitGroupAST(GroupAST *GroupAST) interface{} {
	return nil
}