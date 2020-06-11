package src

import "strings"

// represents a position in code
type Position struct {
	// store a copy of the current line we are processing for reporting errors
	Indent uint32
	Line   uint32
}

// reports messages to the current module
type Reporter struct {
	CurrentLine string
	Position    Position
}

func NewReporter() *Reporter {
	return &Reporter{"", Position{}}
}

// display the current line and position we are processing
func (reporter *Reporter) ReportLine() string {
	// first display the line
	Log(reporter.CurrentLine)
	// then display where we are in that line
	var str strings.Builder
	for i := 0; i < int(reporter.Position.Indent); i++ {
		str.WriteString(" ")
	}
	str.WriteString("^\n")
	for i := 0; i < int(reporter.Position.Indent)-1; i++ {
		str.WriteString("_")
	}
	str.WriteString("/")
	return str.String()
}

func (reporter *Reporter) GetIndent() uint32 {
	return reporter.Position.Indent
}
func (reporter *Reporter) GetLine() uint32 {
	return reporter.Position.Line
}
