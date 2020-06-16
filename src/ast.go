package src

type Visitor interface {
	// statements
	VisitRootAST(RootAST *RootAST) interface{}
	VisitReturnAST(ReturnAST *ReturnAST) interface{}
	VisitBreakAST(BreakAST *BreakAST) interface{}
	VisitForAST(ForAST *ForAST) interface{}
	VisitIfAST(IfAST *IfAST) interface{}
	VisitStructAST(StructAST *StructAST) interface{}
	VisitFnAST(FnAST *FnAST) interface{}
	VisitVarDefAST(VarDefAST *VarDefAST)  interface{}
	VisitBlockAST(BlockAST *BlockAST)  interface{}
	VisitExprSmtAST(ExprStmtAST *ExprStmtAST)  interface{}
	VisitStructSetAST(StructSetAST *StructSetAST) interface{}
	VisitVarSetAST(VarSetAST *VarSetAST) interface{}
	// expressions
	VisitLiteralAST(LiteralAST *LiteralAST)  interface{}
	VisitListAST(ListAST *ListAST)  interface{}
	VisitVariableAST(VariableAST *VariableAST)  interface{}
	VisitUnaryAST(UnaryAST *UnaryAST)  interface{}
	VisitBinaryAST(BinaryAST *BinaryAST)  interface{}
	VisitConnectiveAST(ConnectiveAST *ConnectiveAST)  interface{}
	VisitCallAST(CallAST *CallAST) interface{}
	VisitStructGetAST(StructGet *StructGetAST) interface{}
	VisitGroupAST(GroupAST *GroupAST) interface{}
}

type AST interface {
	Visit(Visitor Visitor) interface{}
}

// statements
type RootAST struct {
	Statements []AST
}
func (RootAST *RootAST) Visit(Visitor Visitor) interface{}{
	return Visitor.VisitRootAST(RootAST)
}
type ReturnAST struct {
	Value AST
}
func (ReturnAST *ReturnAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitReturnAST(ReturnAST)
}
type BreakAST struct{}
func (BreakAST *BreakAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitBreakAST(BreakAST)
}
type ForAST struct {}
func (ForAST *ForAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitForAST(ForAST)
}
type VarSetAST struct {
	Identifier *Token
	Value      AST
}

func (VarSetAST *VarSetAST) Visit(Visitor Visitor) interface{}{
	return Visitor.VisitVarSetAST(VarSetAST)
}

type IfAST struct {
	IfCondition AST
	IfBody      AST

	ElifCondition []AST
	ElifBody      []AST

	ElseBody AST
}
func (IfAST *IfAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitIfAST(IfAST)
}
type StructAST struct {
	Identifier *Token
	Fields     []*VarDefAST
	Packed     bool
}
func (StructAST *StructAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitStructAST(StructAST)
}
type FnAST struct {
	Identifier *Token
	Params     []VarDefAST // the paramaters is an array of definitions
	Body       []AST
	RetType    TavType
	Variadic   bool
}
func (FnAST *FnAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitFnAST(FnAST)
}
type VarDefAST struct {
	Identifier *Token
	Type       TavType
	Assignment AST
}
func (VarDefAST *VarDefAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitVarDefAST(VarDefAST)
}
type BlockAST struct {
	Statements []AST
}
func (BlockAST *BlockAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitBlockAST(BlockAST)
}
type ExprStmtAST struct {
	Expression AST
}
func (ExprStmtAST *ExprStmtAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitExprSmtAST(ExprStmtAST)
}

type LiteralAST struct {
	Type  TavType
	Value TavValue
}
func (LiteralAST *LiteralAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitLiteralAST(LiteralAST)
}
type ListAST struct{}
func (ListAST *ListAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitListAST(ListAST)
}
type VariableAST struct {
	Identifier *Token
}
func (VariableAST *VariableAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitVariableAST(VariableAST)
}
type UnaryAST struct {
	Operator *Token
	Right    AST
}
func (UnaryAST *UnaryAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitUnaryAST(UnaryAST)
}
type BinaryAST struct {
	Left     AST
	Operator *Token
	Right    AST
}
func (BinaryAST *BinaryAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitBinaryAST(BinaryAST)
}
type ConnectiveAST struct {
	Left     AST
	Operator *Token
	Right    AST
}
func (ConnectiveAST *ConnectiveAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitConnectiveAST(ConnectiveAST)
}
type CallAST struct {
	Caller AST
	Args   []AST
}
func (CallAST *CallAST) Visit(Visitor Visitor)  interface{}{
	return Visitor.VisitCallAST(CallAST)
}
type StructGetAST struct {
	Struct AST
	Member *Token
	Deref  bool
}
func (StructGetAST *StructGetAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitStructGetAST(StructGetAST)
}
type StructSetAST struct {
	Struct AST
	Member *Token
	Value  AST
	Deref  bool
}
func (StructSetAST *StructSetAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitStructSetAST(StructSetAST)
}
type GroupAST struct {
	Group AST
}
func (GroupAST *GroupAST) Visit(Visitor Visitor) interface{} {
	return Visitor.VisitGroupAST(GroupAST)
}