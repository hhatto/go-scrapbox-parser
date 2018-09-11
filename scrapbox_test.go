package scrapbox

import (
	"bytes"
	"strings"
	"testing"
)

func execFromString(sb string) string {
	p := NewParser()
	o := p.ToHTML(strings.NewReader(sb))
	outputBuffer := bytes.NewBuffer(o)
	return outputBuffer.String()
}

func TestSimple(t *testing.T) {
	const input = `
aa
 l1
 l2
		l3(use tab*2)
hello. [this is link http://example.com]. cool
hello. [http://example.com this is link] cool
	`
	output := execFromString(input)
	if !strings.Contains(output, "<a href=\"http://example.com\">this is link</a>") {
		t.Errorf("invalid result, href: %v", output)
	}
	if !strings.Contains(output, "<p><span class=\"dot\">l1</span></p>") {
		t.Errorf("invalid result, list1: %v", output)
	}
	if !strings.Contains(output, "<p><span class=\"dot\">l2</span></p>") {
		t.Errorf("invalid result, list2: %v", output)
	}
	if !strings.Contains(output, "<p>aa</p>") {
		t.Errorf("invalid result, simple text: %v", output)
	}
}

func TestHrefTwice(t *testing.T) {
	const input = `
this is [link http://example1.com] and [http://example2.com hello world].
	`
	output := execFromString(input)
	if !strings.Contains(output, "<a href=\"http://example1.com\">link</a>") {
		t.Errorf("invalid result. result: %v", output)
	}
	if !strings.Contains(output, "<a href=\"http://example2.com\">hello world</a>") {
		t.Errorf("invalid result. result: %v", output)
	}
}

func TestStrongTwice(t *testing.T) {
	const input = `
this is [[strong text1]] and [[strong text2]].
	`
	output := execFromString(input)
	if !strings.Contains(output, "<strong>strong text1</strong>") {
		t.Errorf("invalid result. result: %v", output)
	}
	if !strings.Contains(output, "<strong>strong text2</strong>") {
		t.Errorf("invalid result. result: %v", output)
	}
}
