package main

import (
	"log"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// TestParser ... test of parse function
func TestParser(t *testing.T) {
	s := `
	<p>
	hello
	</p>
	 <ol>
	  <li>Dette er en testning. Dette <strong>er</strong> en setning nr. 2</li>
	  <li>Alternativ to
	</ol>
	`
	doc := html.NewTokenizerFragment(strings.NewReader(s), "p")
	var sb strings.Builder
loop:
	for {
		tok := doc.Next()
		log.Println("TOKEN TYPE", tok, doc.Token())
		switch {
		case tok == html.ErrorToken:
			log.Println("ERROR: ", doc.Err())
			break loop
		case tok == html.StartTagToken:
		case tok == html.TextToken:
			sb.Write(doc.Text())
		}
	}
}
