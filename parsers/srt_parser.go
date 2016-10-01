package parser

import (
	"errors"
	"fmt"
	"log"
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
	log.Printf("Parsed number %d", tmp)
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
	for lineNumber := 0; lineNumber < len(lines); lineNumber++ {
		line := lines[lineNumber]
		log.Printf("Line %d: %s", lineNumber, line)
		subtitle := &Subtitle{}
		switch parseState {
		case sNUM:
			// try to parse number of subtitle

			subNumber, err := parseInt(line)
			log.Printf("subNumber %d", subNumber)
			if err != nil {
				return nil, NewParseError("Cannot read subtitle number")
			}
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
			to, err := parseSrtTime(times[0])
			if err != nil {
				return nil, NewParseError("Cannot read `to` time: " + err.Error())
			}
			subtitle.to = to
			subtitle.from = from
			// new state -> subtitle
			parseState = sSUBTITLE
		case sSUBTITLE:
			for {
				// get next
				lineNumber++
				if lineNumber >= len(lines) {
					// save subtitle
					subtitle.text = line
					result = append(result, subtitle)
					break
				}
				// read line
				nextLine := lines[lineNumber]
				// it is the end
				// subtitles with single number after the begining will fail :(
				if _, err := parseInt(nextLine); err != nil {
					// save subtitle
					subtitle.text = line
					result = append(result, subtitle)
					// move backwords
					lineNumber--
					// change state
					parseState = sNUM
					break
				}
				line += "\n" + nextLine
			}
		}
	}
	return result, nil
}
