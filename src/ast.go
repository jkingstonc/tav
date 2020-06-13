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
	Identifier *Token
}

func (StructAST StructAST) Visit() interface{}{
	return nil
}


type FunAST struct {
	Identifier *Token
	Body       *BlockAST
}

func (FunAST FunAST) Visit() interface{}{
	return nil
}

type BlockAST struct {
}

func (BlockAST BlockAST) Visit() interface{}{
	return nil
}