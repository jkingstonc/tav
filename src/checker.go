package src

// implements Visitor
type Checker struct {
	// used for scope checking
	SymTable  *SymTable
}

func Check (compiler *Compiler, RootAST *RootAST) *RootAST{
	return RootAST
}