package scrapbox

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	rgxHrefAfter  = regexp.MustCompile(`\[(.*)\s*(https?://.*)\]`)
	rgxHrefBefore = regexp.MustCompile(`\[(https?://[./a-zA-Z0-9]*)\s(.*)\]`)
	rgxStrong     = regexp.MustCompile(`\[\[([^\]][^\]]*)\]\]`)
	rgxStrongAstr = regexp.MustCompile(`\[\*\s*(.*)\]`)
)

type Parser struct {
	ListLevel int
}

func NewParser() *Parser {
	return &Parser{
		ListLevel: 0,
	}
}

func (p *Parser) ParseList(line string) string {
	match, err := regexp.MatchString(`^\s+`, line)
	if err != nil {
		return line
	}
	if match {
		fmt.Println(match)
		ret := ""
		line = fmt.Sprintf("<span class=\"dot\">%s</span>", ret)
	}
	return line
}

func (p *Parser) ParseHref(line string) string {
	if match := rgxHrefBefore.FindStringSubmatchIndex(line); match != nil {
		before := line[:match[0]]
		after := line[match[1]:]
		link := line[match[2]:match[3]]
		text := line[match[4]:match[5]]
		line = fmt.Sprintf("%s<a href=\"%s\">%s</a>%s", before, link, text, after)
	}
	if match := rgxHrefAfter.FindStringSubmatchIndex(line); match != nil {
		before := line[:match[0]]
		after := line[match[1]:]
		text := strings.TrimSpace(line[match[2]:match[3]])
		link := line[match[4]:match[5]]
		line = fmt.Sprintf("%s<a href=\"%s\">%s</a>%s", before, link, text, after)
	}
	return line
}

func (p *Parser) ParseStrong(line string) string {
	retLine := line
	matchNum := 0
	match := rgxStrong.FindAllStringSubmatchIndex(line, -1)
	if match != nil {
		matchNum = len(match)
	} else {
		return line
	}

	for i := 0; i < matchNum; i++ {
		if m := rgxStrong.FindStringSubmatchIndex(retLine); m != nil {
			retLine = fmt.Sprintf("%s<strong>%s</strong>%s", retLine[:m[0]], retLine[m[2]:m[3]], retLine[m[1]:])
		}
	}
	return retLine
}

func (p *Parser) ToHTML(input io.Reader) []byte {
	var output bytes.Buffer
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		line = p.ParseHref(line)
		line = p.ParseStrong(line)
		output.WriteString(line)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
		return []byte{}
	}
	return output.Bytes()
}