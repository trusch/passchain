package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"
	"github.com/trusch/passchain/cmd/passchain/cmd"
)

func generateManPages() {
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	err := doc.GenManTree(cmd.RootCmd, header, "../../docs/man")
	if err != nil {
		log.Fatal(err)
	}
}

func generateMarkdown() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "../../docs/markdown")
	if err != nil {
		log.Fatal(err)
	}
}

func generateYaml() {
	err := doc.GenYamlTree(cmd.RootCmd, "../../docs/yaml")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	os.MkdirAll("../../docs/man", 0755)
	os.MkdirAll("../../docs/markdown", 0755)
	os.MkdirAll("../../docs/yaml", 0755)
	generateManPages()
	generateMarkdown()
	generateYaml()
}
