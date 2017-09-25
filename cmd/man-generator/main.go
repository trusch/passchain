package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"github.com/trusch/passchain/cmd/passchain/cmd"
)

func main() {

	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	err := doc.GenManTree(cmd.RootCmd, header, "man-pages")
	if err != nil {
		log.Fatal(err)
	}
}
