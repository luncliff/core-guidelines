package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	host     = "raw.githubusercontent.com"
	branch   = "master"
	filename = "CppCoreGuidelines.md"
)

func init() {
	if _, err := os.Stat(filename); os.IsExist(err) {
		return
	}
	source := fmt.Sprintf("https://%s/isocpp/CppCoreGuidelines/%s/%s", host, branch, filename)
	log.Println("source:", source)
	res, err := http.Get(source)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	if err = SaveToFile(res.Body, filename); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func readFile(t *testing.T, filename string) []byte {
	fin, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer fin.Close()
	blob, err := io.ReadAll(fin)
	if err != nil {
		t.Fatal(err)
	}
	return blob
}

func saveBlob(t *testing.T, path string, blob []byte) {
	fout, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer fout.Close()
	if _, err = fout.Write(blob); err != nil {
		t.Fatal(err)
	}
}

func makeNodeSequence(t *testing.T, blob []byte) (sequence []ast.Node) {
	p := goldmark.DefaultParser()
	doc := p.Parse(text.NewReader(blob))
	if doc == nil {
		t.Fail()
	} else if doc.Type() != ast.TypeDocument {
		t.Fatal("parse result is NOT document")
	}
	sequence = make([]ast.Node, 0)
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		sequence = append(sequence, node)
	}
	return
}

func TestMarkdownSectionChunking(t *testing.T) {
	blob := readFile(t, filename)
	if err := os.MkdirAll(path.Join("sections", "en"), 0777); err != nil {
		t.Fatal(err)
	}

	var saveTitle string = "start"
	var pos int = 0
	var nodes []ast.Node = makeNodeSequence(t, blob)
	var headings []*ast.Heading = FilterHeadings(nodes)
	for _, head := range headings {
		segment := head.Lines().At(0)
		title := DropHTML(string(blob[segment.Start:segment.Stop]))
		// // discard some?
		// if strings.Contains(title, "FAQ.") || strings.Contains(title, "Appendix.") {
		// 	continue
		// }
		switch head.Level {
		case 1:
			log.Println(title)
		case 2:
			log.Println("-", title)
		case 3:
			log.Println("  -", title)
		}
		if head.Level < 3 {
			filename := fmt.Sprint(MakeFilename(saveTitle), ".md")
			fpath := path.Join("sections", "en", filename)
			cutidx := segment.Start - head.Level - 1
			saveBlob(t, fpath, blob[pos:cutidx])
			pos = cutidx
			saveTitle = title
		}
	}
}

func mergeSegments(blob []byte, segments *text.Segments) string {
	buf := new(bytes.Buffer)
	for i := 0; i < segments.Len(); i += 1 {
		segment := segments.At(i)
		row := blob[segment.Start:segment.Stop]
		buf.Write(row)
	}
	return buf.String()
}

func TestMarkdownSectionCodeblocks(t *testing.T) {
	blob := readFile(t, filename)

	var nodes []ast.Node = makeNodeSequence(t, blob)
	for _, node := range nodes {
		switch node.Kind() {
		case ast.KindCodeBlock:
			segments := node.Lines()
			if segments == nil {
				log.Println("empty")
				continue
			}
			txt := mergeSegments(blob, segments)
			log.Println("```\n", txt, "\n```")
		}
	}
}
