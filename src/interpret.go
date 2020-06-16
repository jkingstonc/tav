package src

type Interpreter struct {
	Compiler *Compiler
	Reporter *Reporter
	Root     *RootAST
}

func Interpret(Compiler *Compiler, RootAST *RootAST) interface{} {
	return nil
}
