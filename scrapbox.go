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
	rgxSpace      = regexp.MustCompile(`^\s+(\S+)`)
	rgxCode       = regexp.MustCompile(`` + "`" + `([^` + "`" + `]+)` + "`")
)

type Parser struct {
	ListLevel int
}

func NewParser() *Parser {
	return &Parser{
		ListLevel: 0,
	}
}

func (p *Parser) ParseTitle(line string) string {
	return fmt.Sprintf("<h1>%s</h1>", line)
}

func (p *Parser) ParseList(line string) string {
	if line == "" {
		return line
	}
	match := rgxSpace.FindStringSubmatchIndex(line)
	if len(match) > 0 {
		level := match[2]
		line = fmt.Sprintf("<p><span class=\"dot\" style=\"margin-left:%dem;\">🌱&nbsp;%s</span></p>", level, line[match[2]:])
	} else {
		line = fmt.Sprintf("<p>%s</p>", line)
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

func (p *Parser) ParseCode(line string) string {
	retLine := line
	matchNum := 0
	match := rgxCode.FindAllStringSubmatchIndex(line, -1)
	if match != nil {
		matchNum = len(match)
	} else {
		return line
	}

	for i := 0; i < matchNum; i++ {
		if m := rgxCode.FindStringSubmatchIndex(retLine); m != nil {
			retLine = fmt.Sprintf("%s<code>%s</code>%s", retLine[:m[0]], retLine[m[2]:m[3]], retLine[m[1]:])
		}
	}
	return retLine
}

func (p *Parser) ToHTML(input io.Reader) []byte {
	var output bytes.Buffer
	scanner := bufio.NewScanner(input)
	first := true
	for scanner.Scan() {
		line := scanner.Text()
		// inline
		line = p.ParseHref(line)
		line = p.ParseStrong(line)
		line = p.ParseCode(line)

		// block
		if first {
			line = p.ParseTitle(line)
			first = false
		} else {
			line = p.ParseList(line)
		}

		output.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
		return []byte{}
	}
	return output.Bytes()
}
