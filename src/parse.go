package src

type Parser struct {
	Reporter Reporter
	SymTable *SymTable
}

func Parse(compiler *Compiler, tokens []*Token) *AST {
	return nil
}
