package src

import (
	"github.com/llir/llvm/ir/types"
	"os"
)

const (
	WARNING  = 0x0
	CRITICAL = 0x1

	TYPE_SYM_TABLE  uint32 = 0x0
	TYPE_U8         uint32 = 0x1
	TYPE_I8         uint32 = 0x2
	TYPE_U16        uint32 = 0x3
	TYPE_I16        uint32 = 0x4
	TYPE_U32        uint32 = 0x5
	TYPE_I32        uint32 = 0x6
	TYPE_F32        uint32 = 0x7
	TYPE_U64        uint32 = 0x8
	TYPE_I64        uint32 = 0x9
	TYPE_F64        uint32 = 0xA
	TYPE_BOOL       uint32 = 0xB
	TYPE_STRING     uint32 = 0xC
	TYPE_STRUCT     uint32 = 0xD
	TYPE_FN         uint32 = 0xE
	TYPE_ANY        uint32 = 0xF
	TYPE_NULL       uint32 = 0x10
)

type TavValue struct {
	Int    int64
	Float  float64
	String string
	Bool   bool
	Any    interface{}
}

type TavType struct {
	Type    uint32
	IsPtr   bool
	PtrVal  *TavType
	RetType *TavType // used for function calls
}

func (TavType TavType) IsInt() bool{
	return TavType.Type == TYPE_I8 || TavType.Type == TYPE_I16 || TavType.Type == TYPE_I32 || TavType.Type == TYPE_I64
}

func (TavType TavType) IsFloat() bool{
	return TavType.Type == TYPE_F32 || TavType.Type == TYPE_F64
}

type Compiler struct {
	FileName  string
	Source   *string
}

// report an error, the compiler will decide what to do given the severity
func (compiler *Compiler) Report(severity uint8, reporter *Reporter, errCode uint32, msg string) {
	switch severity {
	case WARNING:
		compiler.Warning(reporter, errCode, msg)
	case CRITICAL:
		compiler.Critical(reporter, errCode, msg)
	}
}

// report a warning to the compiler. the compiler will continue and this will not effect the output
func (compiler *Compiler) Warning(reporter *Reporter, errCode uint32, msg string) {
	Log("WARNING", "line", reporter.GetLine(), "pos", reporter.GetIndent(), "in", compiler.FileName,"\n")
	Log(msg)
	Log("\n")
}

// report a critical error to the compiler. the compiler will exit from this point as it cannot continue
func (compiler *Compiler) Critical(reporter *Reporter, errCode uint32, msg string) {
	Log("CRITICAL ERROR", reporter.GetLine(), ":", reporter.GetIndent(), "in", compiler.FileName,"\n")
	Log(reporter.ReportLine())
	Log(msg)
	Log("\n")
	os.Exit(2)
}

func LLType(tavType TavType) types.Type{
	switch tavType.Type {
	case TYPE_BOOL:
		if tavType.IsPtr {
			return types.I1Ptr
		}
		return types.I1
	case TYPE_I8:
		if tavType.IsPtr {
			return types.I8Ptr
		}
		return types.I8
	case TYPE_I16:
		if tavType.IsPtr {
			return types.I16Ptr
		}
		return types.I16
	case TYPE_I32:
		if tavType.IsPtr {
			return types.I32Ptr
		}
		return types.I32
	case TYPE_I64:
		if tavType.IsPtr {
			return types.I64Ptr
		}
		return types.I64
	case TYPE_F32:
		if tavType.IsPtr {
			// NOT SURE HOW THIS WORKS
		}
		return types.Float
	case TYPE_F64:
		if tavType.IsPtr {
			// NOT SURE HOW THIS WORKS
		}
		return types.Double
	}
	return types.Void
}


// TODO some type of check as to whether the inference join was valid
// infer the type of an expression
// this may be somewhat recursive as we have to infer sub types
// if there are multiple types, we do an inference join (infer the correct type given multiple)
func InferType(expression AST, SymTable *SymTable) TavType {
	switch e := expression.(type){
	case *VariableAST:
		t := SymTable.Get(e.Identifier.Lexme())
		Assert(t!=nil, "symbol doesn't exist in symbol table")
		return t.Type
	case *LiteralAST:
		return e.Type
	case *ReturnAST:
		return InferType(e.Value, nil)
	case *BinaryAST:
		return JoinInfered(InferType(e.Left, nil), InferType(e.Right, nil))
	case *CallAST:
		t := InferType(e.Caller, nil)
		return *t.RetType
	}
	// this is unreachable
	return TavType{}
}

// TODO some way to cast the type if they can be joined
// join 2 infered types and figure out what the next type will be
func JoinInfered(type1, type2 TavType) TavType {
	if type1.Type == type2.Type {
		return type1
	}
	// deal with number types
	if type1.IsInt() && type2.IsInt(){
		return type1
	}else if type1.IsFloat() && type2.IsFloat(){
		return type1
	}else if type1.IsFloat() && type2.IsInt(){
		return type1
	}else if type1.IsInt() && type2.IsFloat(){
		return type2
	}
	Assert(false, "cant join infered types")
	return type1
}