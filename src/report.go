package src

import (
	"bufio"
	"strings"
)

// Perhaps the reporter should point to the source string rather than
// holding it in the Reporter struct...

// represents a position in code
type Position struct {
	// store a copy of the current line we are processing for reporting errors
	Indent uint32
	Line   uint32
}

// reports messages to the current module
type Reporter struct {
	Source      *string
	CurrentLine string
	Position    Position
}

func NewReporter(source *string) *Reporter {
	return &Reporter{source, "", Position{}}
}

// display the current line and position we are processing
func (reporter *Reporter) ReportLine() string {
	// get a scanner to find the line that the error is at
	scanner := bufio.NewScanner(strings.NewReader(*reporter.Source))
	for i := 0; i < int(reporter.Position.Line+1); i++ {
		scanner.Scan()
	}
	// first display the line
	Log(scanner.Text())
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

func (reporter *Reporter) GetIndent() uint32 {
	return reporter.Position.Indent
}
func (reporter *Reporter) GetLine() uint32 {
	return reporter.Position.Line
}
