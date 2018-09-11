# Regex based [Scrapbox](https://scrapbox.io) Format Parser

**This is experimental package. Highly under development.**

## Installation

```
$ go get -u -v github.com/hhatto/go-scrapbox-parser
```

## Usage

example of scrapbox format text to HTML

`sb2html.go`:

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	scrapbox "github.com/hhatto/go-scrapbox-parser"
)

func main() {
	filename := os.Args[1]
	p := scrapbox.NewParser()
	input, err := os.Open(filename)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer input.Close()
	o := p.ToHTML(input)
	outputBuffer := bytes.NewBuffer(o)
	output := outputBuffer.String()
	fmt.Println(output)
}
```

`target.scrapbox`:

```
タイトル
[[あ]]
い
 う
 え
  お
 `this is code`.
 [これはリンク http://github.com]です。
[[見出し]]
```

to HTML:

```
$ go run sb2html.go target.scrapbox
<h1>タイトル</h1>
<p><strong>あ</strong></p>
<p>い</p>
<p><span class="dot" style="margin-left:1em;">🌱&nbsp;う</span></p>
<p><span class="dot" style="margin-left:1em;">🌱&nbsp;え</span></p>
<p><span class="dot" style="margin-left:2em;">🌱&nbsp;お</span></p>
<p><span class="dot" style="margin-left:1em;">🌱&nbsp;<code>this is code</code>.</span></p>
<p><span class="dot" style="margin-left:1em;">🌱&nbsp;<a href="http://github.com">これはリンク</a>です。</span></p>
<p><strong>見出し</strong></p>
```

## License
MIT

