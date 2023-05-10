// You can generate the final JSON from the JSON extracted with extract.go
// using:
//   go run clean.go -input raw.json > sabdakosh.json
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	input = flag.String("input", "", "Path to JSON created using 'extract.go'")
)

type rawDefinition struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
}

type definition struct {
	Word    string `json:"word"`
	Defn  []defn `json:"definitions"`
}

var classMap = map[string]string{
	// p preceeds the definition/sense
	"▦": ` class="definition"`,
	// span precedes the initial word
	"▥": ` class="word"`,
	// p precedes section markers for complete definitions (different from
	// variants)
	"▤": ` class="section"`,
	// a part of speech
	"◳": ` class="grammar"`,
	// p followed by an etymology
	"◧": ` class="etymology-marker"`,
	// a etymology
	"◰": ` class="etymology"`,
	// span example
	"▧": ` class="example"`,
}

func classes(in string) string {
	final := in
	for match, replacement := range classMap {
		final = strings.ReplaceAll(final, match, replacement)
	}
	return final
}

func nobreaks(in string) string {
	return strings.ReplaceAll(in, `<br/>`, "")
}

func getClass(n *html.Node) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return attr.Val
		}
	}
	return ""
}

func sections(body *html.Node) [][]*html.Node {
	var final [][]*html.Node
	var cur []*html.Node
	for n := body.FirstChild; n != nil; n = n.NextSibling {
		typ := getClass(n)
		if typ == "section" {
			final = append(final, cur)
			cur = []*html.Node{}
		} else {
			cur = append(cur, n)
		}
	}
	// TODO: Fail if we didn't find any nodes.
	final = append(final, cur)
	return final
}

type defn struct {
	Grammar   string `json:"grammar,omitempty"`
	Etymology string `json:"etymology,omitempty"`
	Senses    []string `json:"senses"`
}

func extract(nodes []*html.Node) defn {
	result := defn{}
	for _, n := range nodes {
		switch getClass(n) {
		// TODO: Case handling needs to be more robust by always extracting text
		// contents instead of relying on tree structure.
		case "etymology-marker":
			result.Etymology = n.FirstChild.FirstChild.Data
		case "grammar":
			result.Grammar = n.FirstChild.Data
		case "definition":
			result.Senses = append(result.Senses, n.FirstChild.Data)
		}
	}
	return result
}

func convert(raw []rawDefinition) []definition {
	elements := []definition{}
	for _, r := range raw {
		cleaned := classes(nobreaks(r.Meaning))
		doc, err := html.ParseFragment(strings.NewReader(cleaned), nil)
		if err != nil {
			log.Fatal(err)
		}
		// From the HTML body (first child is head).
		sections := sections(doc[0].FirstChild.NextSibling)
		defns := []defn{}
		for _, s := range sections {
			d := extract(s)
			defns = append(defns, d)
		}
		elements = append(elements, definition{
			Word:    r.Word,
			Defn: defns,
		})
	}
	return elements
}

func main() {
	flag.Parse()
	data, err := os.ReadFile(*input)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	rds := []rawDefinition{}
	if err := json.Unmarshal(data, &rds); err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}
	ds := convert(rds)
	marshalled, err := json.Marshal(ds)
	if err != nil {
		log.Fatalf("Error marshalling extracted data: %v", err)
	}
	os.Stdout.Write(marshalled)
}
