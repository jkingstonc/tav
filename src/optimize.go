package src

func Optimize(compiler *Compiler, RootAST *RootAST) *RootAST {
	RootAST = Pass1(compiler, RootAST)
	RootAST = Pass2(compiler, RootAST)
	RootAST = Pass3(compiler, RootAST)
	return RootAST
}
