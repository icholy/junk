package fts

import (
	"encoding/xml"
	"os"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

// Document contains a wikipedia extract
type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

// Load documents from xml
func Load(path string) ([]Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var root struct {
		Docs []Document `xml:"doc"`
	}
	dec := xml.NewDecoder(f)
	if err := dec.Decode(&root); err != nil {
		return nil, err
	}
	for i := range root.Docs {
		root.Docs[i].ID = i
	}
	return root.Docs, nil
}

// StopWords should not be indexed or queried
var StopWords = NewSet("a", "and", "be", "have", "i", "in", "or", "that", "the", "to")

// Tokenize splits text into tokens which can be indexed
func Tokenize(text string) []string {
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	r := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		tok = strings.ToLower(tok)
		if !StopWords.Has(tok) {
			r = append(r, english.Stem(tok, false))
		}
	}
	return r
}

// Index is an inverted index
type Index map[string][]int

// Insert populates the index with the documents
func (i Index) Insert(docs []Document) {
	for _, doc := range docs {
		for _, tok := range Tokenize(doc.Text) {
			ids := i[tok]
			if len(ids) == 0 || ids[len(ids)-1] != doc.ID {
				i[tok] = append(ids, doc.ID)
			}
		}
	}
}

// Intersection returns the common integers in two sorted slices
func Intersection(a, b []int) []int {
	maxlen := len(a)
	if n := len(b); n > maxlen {
		maxlen = n
	}
	r := make([]int, 0, maxlen)
	var i, j int
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}

// SearchIDs returns the document IDs which match the query string
func (i Index) SearchIDs(query string) []int {
	var r []int
	for _, tok := range Tokenize(query) {
		if ids, ok := i[tok]; ok {
			if r == nil {
				r = ids
			} else {
				r = Intersection(r, ids)
			}
		}
	}
	return r
}

// Search returns the documents matching query
func (i Index) Search(docs []Document, query string) []Document {
	var r []Document
	for _, id := range i.SearchIDs(query) {
		r = append(r, docs[id])
	}
	return r
}

// Set is a set of strings
type Set map[string]struct{}

// NewSet returns a new set with the provided values
func NewSet(ss ...string) Set {
	set := Set{}
	set.Add(ss...)
	return set
}

// Add values to the set
func (s Set) Add(vv ...string) {
	for _, v := range vv {
		s[v] = struct{}{}
	}
}

// Has returns true if the value is contained in the set
func (s Set) Has(v string) bool {
	_, ok := s[v]
	return ok
}
