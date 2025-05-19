package main

import (
	"log"
	"os"

	"github.com/gkwa/impishimpala/cmd"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmd.Execute()
}
