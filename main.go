package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/jdkato/prose.v2"
)

// Stllinger ... en liste med stillingsannonser.
type Stllinger struct {
	Content []struct {
		UUID          string    `json:"uuid"`
		Published     time.Time `json:"published"`
		Expires       time.Time `json:"expires"`
		WorkLocations []struct {
			Country    string `json:"country"`
			City       string `json:"city"`
			PostalCode string `json:"postalCode"`
			County     string `json:"county"`
			Municipal  string `json:"municipal"`
		} `json:"workLocations"`
		Title          string   `json:"title"`
		Description    string   `json:"description"` // Stillingstekst html5
		Source         string   `json:"source"`
		ApplicationDue string   `json:"applicationDue"`
		Occupations    []string `json:"occupations"`
		Link           string   `json:"link"`
		Employer       struct {
			Name        string `json:"name"`
			Orgnr       string `json:"orgnr"`
			Description string `json:"description"`
		} `json:"employer"`
	} `json:"content"`
	TotalElements int    `json:"totalElements"`
	PageNumber    int    `json:"pageNumber"`
	PageSize      int    `json:"pageSize"`
	TotalPages    int    `json:"totalPages"`
	First         bool   `json:"first"`
	Last          bool   `json:"last"`
	Sort          string `json:"sort"`
}

func main() {

	// Les seneste stillinger med public token.
	bearer := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJwdWJsaWMudG9rZW4udjFAbmF2Lm5vIiwiYXVkIjoiZmVlZC1hcGktdjEiLCJpc3MiOiJuYXYubm8iLCJpYXQiOjE1NTc0NzM0MjJ9.jNGlLUF9HxoHo5JrQNMkweLj_91bgk97ZebLdfx3_UQ"
	url := "https://arbeidsplassen.nav.no/public-feed/api/v1/ads"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("size", "6000")
	q.Add("page", "0")

	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	stllinger := Stllinger{}
	err = json.Unmarshal(body, &stllinger)
	if err != nil {
		log.Fatal(err)
	}

	var sb strings.Builder
	for _, v := range stllinger.Content {
		htmlToString(&sb, v.Description)
	}
	mestBrukteSetninger(sb.String())

}

// Konverterer html5 fragmentene til vanlig tekst.
func htmlToString(sb *strings.Builder, doc string) string {
	// XXX: Trenger å bruke flere formateringshints for å lage setninger.

	d := html.NewTokenizerFragment(strings.NewReader(doc), "p")
loop:
	for {
		tok := d.Next()
		switch {
		case tok == html.ErrorToken:
			break loop
		case tok == html.StartTagToken:
		case tok == html.EndTagToken:
			if d.Token().String() == "</li>" {
				sb.WriteString(".")
			}
		case tok == html.TextToken:
			sb.Write(d.Text())
		}
	}
	return sb.String()
}

// tar tekster og finner repeterte setninger.
func mestBrukteSetninger(texts string) {

	type kv struct {
		Key   string
		Value int
	}

	type Rapport struct {
		AntallSetninger int
		Setningslengde  int
		Toppliste       []kv
	}

	doc, _ := prose.NewDocument(texts, prose.WithSegmentation(true))

	sents := doc.Sentences()

	r := Rapport{}
	r.AntallSetninger = len(sents)

	totlen := 0

	// Set med antall
	var liste = make(map[string]int, 100)
	for _, sentence := range sents {
		_, exists := liste[sentence.Text]
		if !exists {
			if len(strings.TrimSpace(sentence.Text)) > 1 {
				liste[sentence.Text] = 1
				totlen += len(sentence.Text)
			}
		} else {
			liste[sentence.Text]++
		}
	}
	r.Setningslengde = totlen / r.AntallSetninger

	r.Toppliste = make([]kv, 0)
	for k, v := range liste {
		r.Toppliste = append(r.Toppliste, kv{k, v})
	}
	sort.Slice(r.Toppliste, func(i, j int) bool {
		return r.Toppliste[i].Value > r.Toppliste[j].Value
	})
	r.Toppliste = r.Toppliste[0:10]

	template1 := `

	RAPPORT
	=============================================================
	Antall setninger ............. {{.AntallSetninger}}
	Antall ord pr. setning ........{{.Setningslengde}}


	Toppliste
	---------------------------------------------------------------
	Antall
	{{range $k := .Toppliste}}
	{{$k.Value}}	{{ $k.Key}}
	{{end}}

	`

	var out strings.Builder
	templ, _ := template.New("Rapport").Parse(template1)
	templ.Execute(&out, r)
	fmt.Print(out.String())

	// Debug output
	//s, _ := json.MarshalIndent(r, "", "  ")
	//log.Println(string(s))

}
