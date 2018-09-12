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
	rgxHrefAfter    = regexp.MustCompile(`\[(.*)\s*(https?://.*)\]`)
	rgxHrefBefore   = regexp.MustCompile(`\[(https?://[./a-zA-Z0-9]*)\s(.*)\]`)
	rgxRawURL       = regexp.MustCompile(`https?://[\w/:%#\$&\?\(\)~\.=\+\-]+`)
	rgxRawURLIgnore = regexp.MustCompile(`https?://[\w/:%#\$&\?\(\)~\.=\+\-]+\]`)
	//rgxRawURLIgnore = regexp.MustCompile(`[^\[]http(s)?://([\w-]+\.)+[\w-]+(/[\w-./?%&=\]]*)?`)
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
	retLine := line
	if match := rgxHrefBefore.FindAllStringSubmatchIndex(line, -1); match != nil {
		matchNum := len(match)
		for i := 0; i < matchNum; i++ {
			if m := rgxHrefBefore.FindStringSubmatchIndex(retLine); m != nil {
				before := retLine[:m[0]]
				after := retLine[m[1]:]
				link := retLine[m[2]:m[3]]
				text := retLine[m[4]:m[5]]
				retLine = fmt.Sprintf("%s<a href=\"%s\">%s</a>%s", before, link, text, after)
				fmt.Println("G: ", retLine)
			}
		}
	}
	if match := rgxHrefAfter.FindAllStringSubmatchIndex(retLine, -1); match != nil {
		matchNum := len(match)
		for i := 0; i < matchNum; i++ {
			if m := rgxHrefAfter.FindStringSubmatchIndex(retLine); m != nil {
				before := retLine[:m[0]]
				after := retLine[m[1]:]
				text := strings.TrimSpace(retLine[m[2]:m[3]])
				link := retLine[m[4]:m[5]]
				retLine = fmt.Sprintf("%s<a href=\"%s\">%s</a>%s", before, link, text, after)
			}
		}
	}
	return retLine
}

func (p *Parser) ParseRawURL(line string) string {
	fmt.Println("line: ", line)
	retLine := ""
	checkLine := line
	matchNum := 0
	match := rgxRawURL.FindAllStringSubmatchIndex(line, -1)
	if match != nil {
		matchNum = len(match)
		fmt.Println("raw-match:", matchNum, match, line[match[0][0]:match[0][1]])
	} else {
		return line
	}

	offset := 0
	for i := 0; i < matchNum; i++ {
		if m := rgxRawURL.FindStringSubmatchIndex(checkLine); m != nil {
			link := checkLine[m[0]:m[1]]
			fmt.Printf("VV0: link='%s', check='%s'\n", link, checkLine)
			fmt.Printf("VV0: check[='%c', check]='%c'\n", checkLine[m[0]-1], checkLine[m[1]])
			if checkLine[m[0]-1] == '[' {
				fmt.Printf("[-ignore: %v\n", m)
				offset = m[1]
				retLine += checkLine[:m[1]]
				checkLine = checkLine[m[1]:]
				continue
			} else if checkLine[m[1]] == ']' {
				fmt.Printf("]-ignore: %v\n", m)
				offset = m[1]
				retLine += checkLine[:m[1]]
				checkLine = checkLine[m[1]:]
				continue
			}
			//if mIgnore := rgxRawURLIgnore.FindStringSubmatchIndex(checkLine); mIgnore != nil {
			//	ignoreLink := checkLine[mIgnore[0]:mIgnore[1]]
			//	fmt.Printf("ignore: '%v', '%s'\n", mIgnore, ignoreLink)
			//	if string(ignoreLink[len(ignoreLink)-1]) == "]" {
			//		fmt.Println("IGNORE:", mIgnore, ignoreLink)
			//		retLine += checkLine[:mIgnore[1]]
			//		checkLine = checkLine[mIgnore[1]:]
			//		continue
			//	}
			//}
			retLine += fmt.Sprintf("%s<a href=\"%s\">%s</a>%s", checkLine[:m[0]], link, link, checkLine[m[1]:])
			fmt.Println("vv:", retLine)
		}
	}
	if offset != 0 {
		fmt.Println(offset, checkLine)
		retLine += checkLine[:]
	}
	fmt.Println("RET:", retLine)
	return retLine
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
		line = p.ParseRawURL(line)
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
