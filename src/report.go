package src

import "strings"

// represents a position in code
type Position struct {
	// store a copy of the current line we are processing for reporting errors
	Indent uint32
	Line   uint32
}

type Reporter interface {
	ReportLine() string
	GetIndent() uint32
	GetLine() uint32
}

// reports messages to the current module
type LexReporter struct {
	CurrentLine string
	Position    Position
}

// display the current line and position we are processing
func (reporter LexReporter) ReportLine() string {
	// first display the line
	Log(reporter.CurrentLine)
	// then display where we are in that line
	var str strings.Builder
	for i := 0; i < int(reporter.Position.Indent)-1; i++ {
		str.WriteString(" ")
	}
	str.WriteString("^\n")
	for i := 0; i < int(reporter.Position.Indent)-2; i++ {
		str.WriteString("_")
	}
	str.WriteString("/")
	return str.String()
}

func (reporter LexReporter) GetIndent() uint32 {
	return reporter.Position.Indent
}
func (reporter LexReporter) GetLine() uint32 {
	return reporter.Position.Line
}
