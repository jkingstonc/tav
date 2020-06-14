package src

type Visitor interface {
	// statements
	VisitRootAST(RootAST *RootAST)
	VisitReturnAST(ReturnAST *ReturnAST)
	VisitBreakAST(BreakAST *BreakAST)
	VisitForAST(ForAST *ForAST)
	VisitIfAST(IfAST *IfAST)
	VisitStructAST(StructAST *StructAST)
	VisitFnAST(FnAST *FnAST)
	VisitVarDefAST(VarDefAST *VarDefAST)
	VisitBlockAST(BlockAST *BlockAST)
	VisitExprSmtAST(ExprStmtAST *ExprStmtAST)
	// expressions
	VisitAssignAST(AsssignAST *AsssignAST)
	VisitLiteralAST(LiteralAST *LiteralAST)
	VisitListAST(ListAST *ListAST)
	VisitVariableAST(VariableAST *VariableAST)
	VisitUnaryAST(UnaryAST *UnaryAST)
	VisitBinaryAST(BinaryAST *BinaryAST)
	VisitConnectiveAST(ConnectiveAST *ConnectiveAST)
	VisitCallAST(CallAST *CallAST)
	VisitStructGetAST(StructGet *StructGetAST)
	VisitStructSetAST(StructSetAST *StructSetAST)
	VisitGroupAST(GroupAST *GroupAST)
}

type AST interface {
	Visit(Visitor Visitor)
}

// statements
type RootAST struct {
	Statements []AST
}
func (RootAST *RootAST) Visit(Visitor Visitor) {
	Visitor.VisitRootAST(RootAST)
}
type ReturnAST struct {
	Value AST
}
func (ReturnAST *ReturnAST) Visit(Visitor Visitor) {
	Visitor.VisitReturnAST(ReturnAST)
}
type BreakAST struct{}
func (BreakAST *BreakAST) Visit(Visitor Visitor) {
	Visitor.VisitBreakAST(BreakAST)
}
type ForAST struct {}
func (ForAST *ForAST) Visit(Visitor Visitor) {
	Visitor.VisitForAST(ForAST)
}
type IfAST struct {
	IfCondition AST
	IfBody      AST

	ElifCondition []AST
	ElifBody      []AST

	ElseBody AST
}
func (IfAST *IfAST) Visit(Visitor Visitor) {
	Visitor.VisitIfAST(IfAST)
}
type StructAST struct {
	Identifier *Token
}
func (StructAST *StructAST) Visit(Visitor Visitor) {
	Visitor.VisitStructAST(StructAST)
}
type FnAST struct {
	Identifier *Token
	Params     []VarDefAST // the paramaters is an array of definitions
	Body       *BlockAST
	RetType    TavType
	Variadic   bool
}
func (FnAST *FnAST) Visit(Visitor Visitor) {
	Visitor.VisitFnAST(FnAST)
}
type VarDefAST struct {
	Identifier *Token
	Type       TavType
	Assignment AST
}
func (VarDefAST *VarDefAST) Visit(Visitor Visitor) {
	Visitor.VisitVarDefAST(VarDefAST)
}
type BlockAST struct {
	Statements []AST
}
func (BlockAST *BlockAST) Visit(Visitor Visitor) {
	Visitor.VisitBlockAST(BlockAST)
}
type ExprStmtAST struct {
	Expression AST
}
func (ExprStmtAST *ExprStmtAST) Visit(Visitor Visitor) {
	Visitor.VisitExprSmtAST(ExprStmtAST)
}

// expressions
type AsssignAST struct { // assign to a non-member variable
	Identifier *Token
	Value      AST
}
func (AsssignAST *AsssignAST) Visit(Visitor Visitor) {
	Visitor.VisitAssignAST(AsssignAST)
}
type LiteralAST struct {
	Type  uint32
	Value interface{}
}
func (LiteralAST *LiteralAST) Visit(Visitor Visitor) {
	Visitor.VisitLiteralAST(LiteralAST)
}
type ListAST struct{}
func (ListAST *ListAST) Visit(Visitor Visitor) {
	Visitor.VisitListAST(ListAST)
}
type VariableAST struct {
	Identifier *Token
}
func (VariableAST *VariableAST) Visit(Visitor Visitor) {
	Visitor.VisitVariableAST(VariableAST)
}
type UnaryAST struct {
	Operator *Token
	Right    AST
}
func (UnaryAST *UnaryAST) Visit(Visitor Visitor) {
	Visitor.VisitUnaryAST(UnaryAST)
}
type BinaryAST struct {
	Left     AST
	Operator *Token
	Right    AST
}
func (BinaryAST *BinaryAST) Visit(Visitor Visitor) {
	Visitor.VisitBinaryAST(BinaryAST)
}
type ConnectiveAST struct {
	Left     AST
	Operator *Token
	Right    AST
}
func (ConnectiveAST *ConnectiveAST) Visit(Visitor Visitor) {
	Visitor.VisitConnectiveAST(ConnectiveAST)
}
type CallAST struct {
	Caller AST
	Args   []AST
}
func (CallAST *CallAST) Visit(Visitor Visitor) {
	Visitor.VisitCallAST(CallAST)
}
type StructGetAST struct {
	Struct AST
	Member *Token
}
func (StructGetAST *StructGetAST) Visit(Visitor Visitor) {
	Visitor.VisitStructGetAST(StructGetAST)
}
type StructSetAST struct {
	Struct AST
	Member *Token
	Value  AST
}
func (StructSetAST *StructSetAST) Visit(Visitor Visitor) {
	Visitor.VisitStructSetAST(StructSetAST)
}
type GroupAST struct {
	Group AST
}
func (GroupAST *GroupAST) Visit(Visitor Visitor) {
	Visitor.VisitGroupAST(GroupAST)
}