package src


// Compile to LLVM and JIT compile
func JIT(compiler *Compiler, ast *AST) (uint8, string) {
	LLVMCompile(compiler, ast)
	return SUCCESS_COMP, "lol"
}