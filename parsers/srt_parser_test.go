package parser

import "testing"

func TestOneSubtitle(t *testing.T) {
	srtParser := &SrtParser{}
	oneSubtitle := `1
00:00:12,014 --> 00:00:14,517
[theme music playing]`
	result, err := srtParser.Parse(oneSubtitle)
	if err != nil {
		t.Errorf("No error expected %v", err)
		return
	}
	if result == nil {
		t.Error("Result expected")
		return
	}
	if len(result) != 1 {
		t.Errorf("One subtitle expected, got %d", len(result))
		return
	}
	subtitle := result[0]
	t.Logf("subtitle %v", subtitle)
	if subtitle.num != 1 {
		t.Errorf("Subtitle number parsed wrong, got %v", subtitle.num)
		return
	}
	if subtitle.from != 12014 {
		t.Errorf("From parsed wrong, got %v", subtitle.from)
		return
	}
	if subtitle.from != 14517 {
		t.Errorf("To parsed wrong, got %v", subtitle.to)
		return
	}
	if subtitle.text != "[theme music playing]" {
		t.Errorf("Text extracted wrong")
		return
	}
}
