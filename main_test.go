package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// TestParser ... test of parse function, just a spin.
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
		switch {
		case tok == html.ErrorToken:
			break loop
		case tok == html.StartTagToken:
		case tok == html.TextToken:
			sb.Write(doc.Text())
		}
	}
}
