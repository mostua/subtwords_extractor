package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SrtParser parser compatible with srt format
type SrtParser struct {
}

func parseInt(val string) (int, error) {
	tmp, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(tmp), nil
}

// returns time in miliseconds
func parseSrtTime(time string) (result int, err error) {
	// split by , to get milis
	times := strings.Split(time, ",")
	if len(times) != 2 {
		return 0, errors.New("wrong time format expected xx:xx:xx,xxx")
	}
	// afet ,
	milis, err := parseInt(times[1])
	if err != nil {
		return 0, errors.New("cannot parse milis")
	}
	// before ,
	rest := strings.Split(times[0], ":")
	seconds := 0
	// hours:minutes:seconds
	for i := 0; i < len(rest); i++ {
		seconds *= 60
		val, err := parseInt(rest[i])
		if err != nil {
			return 0, err
		}
		seconds += val
	}
	return seconds*1000 + milis, nil
}

const (
	sNUM      = iota
	sTIME     = iota
	sSUBTITLE = iota
)

// Parse implements parsing of srt subtitles
func (parser *SrtParser) Parse(subitlesTxt string) (result []*Subtitle, parseError *ParsingError) {
	// inital state
	parseState := sNUM
	// lets try with normalize windows
	subitlesTxt = strings.Replace(subitlesTxt, "\r\n", "\n", -1)
	// and then split in the linux way
	lines := strings.Split(subitlesTxt, "\n")
	var subtitle *Subtitle
	for lineNumber := 0; lineNumber < len(lines); lineNumber++ {
		line := lines[lineNumber]
		// omit empty lines
		if line == "" {
			continue
		}
		switch parseState {
		case sNUM:
			// try to parse number of subtitle

			subNumber, err := parseInt(line)
			if err != nil {
				return nil, NewParseError("Cannot read subtitle number")
			}
			subtitle = &Subtitle{}
			subtitle.num = subNumber
			// new state -> time
			parseState = sTIME
		case sTIME:
			times := strings.Split(line, " --> ")
			if len(times) != 2 {
				return nil, NewParseError(fmt.Sprintf("From and to times expected, got: %v", times))
			}
			from, err := parseSrtTime(times[0])
			if err != nil {
				return nil, NewParseError("Cannot read `from` time: " + err.Error())
			}
			to, err := parseSrtTime(times[1])
			if err != nil {
				return nil, NewParseError("Cannot read `to` time: " + err.Error())
			}
			subtitle.to = to
			subtitle.from = from
			// new state -> subtitle
			parseState = sSUBTITLE
		case sSUBTITLE:
			// move to nextLine
			subtitle.text = line
			lineNumber++
			for lineNumber < len(lines) {
				// read line
				nextLine := lines[lineNumber]
				// it is the end ?
				if nextLine == "" {
					break
				}
				subtitle.text += "\n" + nextLine
				lineNumber++
			}
			result = append(result, subtitle)
			parseState = sNUM
		}
	}
	return result, nil
}
