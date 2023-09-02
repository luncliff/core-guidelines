package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"testing"
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

func TestMarkdownChunking(t *testing.T) {
	blob, err := ReadFromFile(filename)
	if err != nil {
		t.Fatal(err)
	}

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
	blob, err := ReadFromFile(filename)
	if err != nil {
		t.Fatal(err)
	}

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
	blob, err := ReadFromFile(filename)
	if err != nil {
		t.Fatal(err)
	}

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
