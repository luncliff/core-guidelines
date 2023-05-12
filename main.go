package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/yuin/goldmark/ast"
)

var (
	source  string
	verbose bool
)

func main() {
	command := os.Args[1]
	switch command {
	case "sync":
		flag.StringVar(&source, "source", "https://github.com/<org>/<repo>", "URL to download guideline markdown file")
	case "help":
		break
	}
	flag.BoolVar(&verbose, "verbose", false, "Print more log messages")

	os.Args = os.Args[1:]
	if flag.Parse(); !flag.Parsed() {
		log.Fatalln("os.Args parse failed")
		os.Exit(1)
	}

	log.Println("source:", source)
	log.Println("verbose:", verbose)
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

func FilterHeadings(nodes []ast.Node) (headings []*ast.Heading) {
	headings = make([]*ast.Heading, 0)
	for _, node := range nodes {
		if node.Type() != ast.TypeBlock {
			continue
		}
		switch node.Kind() {
		case ast.KindHeading:
			var p interface{} = node
			head := p.(*ast.Heading)
			headings = append(headings, head)
		}
	}
	return
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
