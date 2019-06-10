package main

import "github.com/jdkato/prose/summarize"

// Document ... Alias for Document.
type Document summarize.Document

// WordDensityLang ... returns a map of each word and its density.
func (d *Document) WordDensityLang(lang string) map[string]float64 {
	density := make(map[string]float64)

	return density
}
