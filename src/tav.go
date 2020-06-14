package src

import (
	"github.com/llir/llvm/ir/types"
	"os"
)

const (
	WARNING  = 0x0
	CRITICAL = 0x1
)

type TavType struct {
	Type    uint32
	IsPtr   bool
	PtrVal  *TavType
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
	case I8:
		if tavType.IsPtr {
			return types.I8Ptr
		}
		return types.I8
	case I16:
		if tavType.IsPtr {
			return types.I16Ptr
		}
		return types.I16
	case I32:
		if tavType.IsPtr {
			return types.I32Ptr
		}
		return types.I32
	case I64:
		if tavType.IsPtr {
			return types.I64Ptr
		}
		return types.I64
	}
	return types.Void
}
