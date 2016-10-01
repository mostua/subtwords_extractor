package parser

import "testing"

type tWrapper struct {
	*testing.T
}

func (t *tWrapper) validateParsing(result []*Subtitle, err *ParsingError, expectedLength int) bool {
	if err != nil {
		t.Errorf("No error expected %v", err)
		return false
	}
	if result == nil {
		t.Error("Result expected")
		return false
	}
	if len(result) != expectedLength {
		t.Errorf("One subtitle expected, got %d", len(result))
		return false
	}
	return true
}

func (t *tWrapper) validateSub(tested *Subtitle, expected *Subtitle) bool {
	if tested.num != expected.num {
		t.Errorf("Subtitle number parsed wrong, got %v expected %v", tested.num, expected.num)
		return false
	}
	if tested.from != expected.from {
		t.Errorf("From parsed wrong, got %v expected %v", tested.from, expected.from)
		return false
	}
	if tested.to != expected.to {
		t.Errorf("To parsed wrong, got %v expected %v", tested.to, expected.to)
		return false
	}
	if tested.text != expected.text {
		t.Errorf("Text extracted wrong, got %v expected %v", tested.text, expected.text)
		return false
	}
	return true
}

func TestOneSubtitle(t *testing.T) {
	srtParser := &SrtParser{}
	oneSubtitle := `1
00:00:12,014 --> 00:00:14,517
[theme music playing]
`
	result, err := srtParser.Parse(oneSubtitle)
	wrapper := &tWrapper{t}
	if !wrapper.validateParsing(result, err, 1) {
		return
	}
	subtitle := result[0]
	t.Logf("subtitle %v", subtitle)
	if !wrapper.validateSub(subtitle, &Subtitle{1, 12014, 14517, "[theme music playing]"}) {
		return
	}
}

func TestTwoSubtitle(t *testing.T) {
	srtParser := &SrtParser{}
	twoSubtitles := `1
00:00:12,014 --> 00:00:14,517
[theme music playing]

2
00:01:11,111 --> 00:02:22,222
Testing multiline
line

`
	result, err := srtParser.Parse(twoSubtitles)
	wrapper := &tWrapper{t}
	if !wrapper.validateParsing(result, err, 2) {
		return
	}
	if !wrapper.validateSub(result[0], &Subtitle{1, 12014, 14517, "[theme music playing]"}) {
		return
	}

	if !wrapper.validateSub(result[1], &Subtitle{2, 71111, 142222, `Testing multiline
line`}) {
		return
	}
}
