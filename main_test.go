package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"testing"
)

const (
	host     = "raw.githubusercontent.com"
	filename = "CppCoreGuidelines.md"
)

func init() {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		source := fmt.Sprintf("https://%s/isocpp/CppCoreGuidelines/%s/%s", host, "master", filename)
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

func TestMarkdownChunking(t *testing.T) {
	blob := readFile(t, filename)

	folder := path.Join("sections", "en")
	if err := os.RemoveAll(folder); err != nil {
		t.Fatal(err)
	}
	t.Log("removed", folder)
	if err := os.MkdirAll(folder, 0777); err != nil {
		t.Fatal(err)
	}

	nodes, err := MakeNodeSequence(blob)
	if err != nil {
		t.Fatal(err)
	}

	if err := SaveSections(blob, nodes, folder); err != nil {
		t.Fatal(err)
	}
}

func TestMarkdownDecorateCodeBlocks(t *testing.T) {
	blob := readFile(t, filename)

	nodes, err := MakeNodeSequence(blob)
	if err != nil {
		t.Fatal(err)
	}

	fout, err := os.Create("cpp-1.md")
	if err != nil {
		t.Fatal(err)
	}
	defer fout.Close()

	if err = DecorateCodeBlocks(blob, nodes, fout); err != nil {
		t.Fatal(err)
	}
}

func TestMarkdownH5Examples(t *testing.T) {
	blob := readFile(t, filename)

	nodes, err := MakeNodeSequence(blob)
	if err != nil {
		t.Fatal(err)
	}

	fout, err := os.Create("cpp-2.md")
	if err != nil {
		t.Fatal(err)
	}
	defer fout.Close()

	if err = DecorateH5Examples(blob, nodes, fout); err != nil {
		t.Fatal(err)
	}
}
