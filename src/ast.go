package src

type Visitor interface {
	VisitStruct(s *StructAST) interface{}
	VisitFun(f *FunAST) interface{}
}

type AST interface {
	Visit() interface{}
}

type RootAST struct {
	Statements    []AST
}

func (RootAST RootAST) Visit() interface{} {
	return nil
}

type StructAST struct {
}

func (StructAST StructAST) Visit() interface{}{
	return nil
}


type FunAST struct {
	Body     *BlockAST
	RetType  uint32
}

func (FunAST FunAST) Visit() interface{}{
	return nil
}

type DefineAST struct {
	Identifier *Token
	Type 	   uint32
	Assignment AST
}

func (DefineAST DefineAST) Visit() interface{}{
	return nil
}

// an expression e.g. (1+2) or (x.Buffer = 1)
type ExprAST struct{

}

func (ExprAST ExprAST) Visit() interface{}{
	return nil
}
type BlockAST struct {
}

func (BlockAST BlockAST) Visit() interface{}{
	return nil
}