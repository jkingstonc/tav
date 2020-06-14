package src

import "os"

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
