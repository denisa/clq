package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/denisa/clq/internal/validator"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func main() {
	os.Exit(entryPoint(os.Args[0], os.Args[1:]...))
}

func entryPoint(name string, arguments ...string) int {
	options := flag.NewFlagSet(name, flag.ContinueOnError)
	options.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s { flags } <path to changelog.md>\nOptions are:\n", options.Name())
		options.PrintDefaults()
	}
	var release = options.Bool("release", false, "Enable release-mode validation")
	if options.Parse(arguments) != nil {
		return 2
	}

	var documents []string
	if options.NArg() == 0 {
		documents = []string{"-"}
	} else {
		documents = options.Args()
	}
	for _, document := range documents {
		data, err := readInput(document)
		if err != nil {
			log.Println(err)
			return 1
		}

		validator := validator.NewRenderer(validator.WithRelease(*release))
		md := goldmark.New(
			goldmark.WithRenderer(
				renderer.NewRenderer(renderer.WithNodeRenderers(
					util.Prioritized(validator, 1000)))),
		)

		fmt.Print("Processing ", document, "...")
		var buf bytes.Buffer
		if err := md.Convert(data, &buf); err != nil {
			fmt.Println("❗️")
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Println("✅")

	}
	return 0
}

func readInput(input string) ([]byte, error) {
	if input == "-" {
		return ioutil.ReadAll(os.Stdin)
	}

	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
