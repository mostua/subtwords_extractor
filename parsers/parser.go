package parser

// Subtitle Definition of subtitle
type Subtitle struct {
	num  int
	from int
	to   int
	text string
}

// ParsingError represents error of subtitles parsing
type ParsingError struct {
	ErrorMsg string
}

// NewParseError craetes parsing error
func NewParseError(message string) *ParsingError {
	result := &ParsingError{}
	result.ErrorMsg = message
	return result
}

// Parser What is a parser
type Parser interface {
	Parse(string) ([]*Subtitle, *ParsingError)
}
