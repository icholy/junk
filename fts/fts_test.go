package fts

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func load(t *testing.T) []Document {
	t.Helper()
	home, err := os.UserHomeDir()
	assert.NilError(t, err)
	name := filepath.Join(home, "Downloads", "enwiki-latest-abstract1.xml")
	docs, err := Load(name)
	assert.NilError(t, err)
	return docs
}

func TestSearch(t *testing.T) {
	t.Log("loading ...")
	docs := load(t)

	t.Log("indexing ...")
	idx := Index{}
	idx.Insert(docs)

	t.Log("searching ...")
	for _, doc := range idx.Search(docs, "Small wild cat") {
		fmt.Println(doc.ID, doc.Title)
	}
}
