package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	host     = "raw.githubusercontent.com"
	filename = "CppCoreGuidelines.md"
)

func main() {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		source := fmt.Sprintf("https://%s/isocpp/CppCoreGuidelines/%s/%s", host, "master", filename)
		log.Println("source:", source)
		res, err := http.Get(source)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if err = SaveToFile(res.Body, filename); err != nil {
			log.Fatalln(err)
			return
		}
	}
	folder := path.Join("sections", "en")
	if err := os.MkdirAll(folder, 0777); err != nil {
		log.Fatalln(err)
		return
	}

	rewrite1 := func(filename string) error {
		blob, err := ReadFromFile(filename)
		if err != nil {
			return err
		}
		nodes, err := MakeNodeSequence(blob)
		if err != nil {
			return err
		}
		fout, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer fout.Close()
		return DecorateH5Examples(blob, nodes, fout)
	}
	if err := rewrite1(filename); err != nil {
		log.Fatalln(err)
		return
	}

	rewrite2 := func(filename string) error {
		blob, err := ReadFromFile(filename)
		if err != nil {
			return err
		}
		nodes, err := MakeNodeSequence(blob)
		if err != nil {
			return err
		}
		fout, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer fout.Close()
		return DecorateCodeBlocks(blob, nodes, fout)
	}
	if err := rewrite2(filename); err != nil {
		log.Fatalln(err)
		return
	}

	splitToFolder := func(filename, folder string) error {
		blob, err := ReadFromFile(filename)
		if err != nil {
			return err
		}
		nodes, err := MakeNodeSequence(blob)
		if err != nil {
			return err
		}
		return SaveSections(blob, nodes, folder)
	}

	if err := splitToFolder(filename, folder); err != nil {
		log.Fatalln(err)
		return
	}
}

// Discard <a></a> HTML tag
func DropHTML(txt string) string {
	idx := strings.Index(txt, "</a>")
	if idx == -1 {
		return txt
	}
	return txt[idx+4:]
}

// replace some bad characters so the output string can be used as files' name
func MakeFilename(input string) string {
	const t0 = ""
	const t1 = "_"
	const t2 = "-"
	r := strings.NewReplacer( //
		".", t2, ":", t1, ",", t0, //
		"(", t0, ")", t0, //
		" ", t2, "/", t2, //
		"\"", t0, "'", t0)
	return r.Replace(input)
}

func MakeDocument(blob []byte) (ast.Node, error) {
	p := goldmark.DefaultParser()
	doc := p.Parse(text.NewReader(blob))
	if doc == nil {
		return nil, errors.New("failed to parse document")
	} else if doc.Type() != ast.TypeDocument {
		return nil, errors.New("parse result is not a document")
	}
	return doc, nil
}

func MakeNodeSequence(blob []byte) ([]ast.Node, error) {
	doc, err := MakeDocument(blob)
	if err != nil {
		return nil, err
	}
	sequence := make([]ast.Node, 0)
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		sequence = append(sequence, node)
	}
	return sequence, nil
}

func FilterHeadings(nodes []ast.Node, headings chan<- *ast.Heading, maxLevel int) {
	defer close(headings)
	for _, node := range nodes {
		if node.Type() != ast.TypeBlock {
			continue
		}
		switch node.Kind() {
		case ast.KindHeading:
			var p interface{} = node
			heading := p.(*ast.Heading)
			if heading.Level <= maxLevel {
				headings <- heading
			}
		}
	}
}

func GetShortTitle(title string) (string, bool) {
	idx := strings.Index(title, ": ")
	if idx == -1 {
		return title, false
	}
	return title[:idx], true
}

func SaveSections(fullText []byte, nodes []ast.Node, folder string) error {
	headings := make(chan *ast.Heading)

	// When the end of "Heading 1" or "Heading 2" reached,
	// slice the raw text 'pos' and save to the file
	go FilterHeadings(nodes, headings, 2)
	var pos int = 0
	var saveTitle string = "empty"

	for head := range headings {
		segment := head.Lines().At(0)
		title := DropHTML(string(fullText[segment.Start:segment.Stop]))
		filename := fmt.Sprint(MakeFilename(saveTitle), ".md")
		if shorten, found := GetShortTitle(saveTitle); found {
			log.Println(shorten, ":", path.Join(folder, filename))
		}

		cutidx := segment.Start - head.Level - 1
		if err := SaveBlobToFile(fullText[pos:cutidx], path.Join(folder, filename)); err != nil {
			log.Println(err)
			return err
		}
		pos = cutidx
		saveTitle = title
	}
	return nil
}

func ReadFromFile(filename string) ([]byte, error) {
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	return io.ReadAll(fin)
}

func SaveToFile(rc io.ReadCloser, filename string) error {
	defer rc.Close()
	fout, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fout.Close()
	_, err = io.Copy(fout, rc)
	return err
}

func SaveBlobToFile(blob []byte, path string) error {
	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fout.Close()
	_, err = fout.Write(blob)
	return err
}

func MergeSegments(blob []byte, segments *text.Segments) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < segments.Len(); i += 1 {
		segment := segments.At(i)
		row := blob[segment.Start:segment.Stop]
		buf.Write(row)
	}
	return buf.Bytes()
}

func DecorateCodeBlocks(blob []byte, nodes []ast.Node, writer io.Writer) error {
	var pos int = 0
	for _, node := range nodes {
		segments := node.Lines()
		switch node.Kind() {
		case ast.KindCodeBlock: // if code block, mark as C++ code
			if segments == nil {
				log.Println("no segment")
				continue
			}
			block := MergeSegments(blob, segments)
			writer.Write([]byte("\n\n```c++\n"))
			writer.Write(block)
			writer.Write([]byte("```\n"))
			segment := segments.At(segments.Len() - 1)
			pos = segment.Stop
			continue
		}
		// bypass the other segments (and some gaps between previous segments)
		for i := 0; i < segments.Len(); i += 1 {
			segment := segments.At(i)
			row := blob[pos:segment.Stop]
			writer.Write(row)
			pos = segment.Stop
		}
	}
	return nil
}

func MakeAdmonition(text string) (value string, indent bool) {
	guessAdmonition := func() (string, bool) {
		if strings.Contains(text, "Reason") {
			return "info", true
		}
		if strings.Contains(text, "Example") {
			if strings.Contains(text, ", good") {
				return "success", false
			}
			if strings.Contains(text, ", bad") {
				return "failure", false
			}
			return "example", false
		}
		if strings.Contains(text, "Enforcement") {
			return "tip", true
		}
		if strings.Contains(text, "Discussion") {
			return "quote", true
		}
		if strings.Contains(text, "Exception") {
			return "warning", true
		}
		// "See", "See also", "Alternative" "Alternatives"
		return "note", false
	}
	value, indent = guessAdmonition()
	value = fmt.Sprintf("!!! %s \"%s\"", value, text)
	return
}

func DecorateH5Examples(blob []byte, nodes []ast.Node, writer io.Writer) error {
	var pos int = 0
	for _, node := range nodes {
		segments := node.Lines()
		switch node.Kind() {
		case ast.KindHeading:
			heading := node.(*ast.Heading)
			// mostly 5, but sometimes 4 (probably mistyped)
			if heading.Level >= 4 {
				segment := segments.At(0)
				cutidx := segment.Start - heading.Level - 1
				writer.Write(blob[pos:cutidx])
				adomination, _ := MakeAdmonition(string(blob[segment.Start:segment.Stop]))
				writer.Write([]byte(adomination))
				pos = segment.Stop
				continue
			}
		}
		// bypass the other segments (and some gaps between previous segments)
		for i := 0; i < segments.Len(); i += 1 {
			segment := segments.At(i)
			row := blob[pos:segment.Stop]
			writer.Write(row)
			pos = segment.Stop
		}
	}
	return nil
}
