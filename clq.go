package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/denisa/clq/internal/query"
	"github.com/denisa/clq/internal/validator"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Clq struct {
	stdin          io.Reader
	stdout, stderr io.Writer
}

func main() {
	clq := &Clq{stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr}
	os.Exit(clq.entryPoint(os.Args[0], os.Args[1:]...))
}

func (clq *Clq) entryPoint(name string, arguments ...string) int {
	options := flag.NewFlagSet(name, flag.ContinueOnError)
	options.SetOutput(clq.stderr)
	options.Usage = func() {
		fmt.Fprintf(options.Output(), "\nUsage: %s { flags } <path to changelog.md>\n\nOptions are:\n", options.Name())
		options.PrintDefaults()
	}
	var release = options.Bool("release", false, "Enable release-mode validation")
	var queryString = options.String("query", "", "A query to extract information out of the change log")
	if options.Parse(arguments) != nil {
		return 2
	}

	var documents []string
	if options.NArg() == 0 {
		documents = []string{"-"}
	} else {
		documents = options.Args()
	}

	var hasError bool
	for _, document := range documents {
		queryEngine, err := query.NewQueryEngine(*queryString)
		if err != nil {
			fmt.Fprintf(clq.stderr, "❗️ %v\n", err)
			hasError = true
			return 2
		}

		source, err := clq.readInput(document)
		if err != nil {
			if len(documents) == 1 {
				fmt.Fprintf(clq.stderr, "❗️ %v\n", err)
			} else {
				fmt.Fprintf(clq.stderr, "❗️ %v: %v\n", document, err)
			}
			hasError = true
			continue
		}

		reader := text.NewReader(source)
		doc := goldmark.DefaultParser().Parse(reader)

		validationEngine := renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(
					validator.NewRenderer(
						validator.WithQuery(*queryEngine),
						validator.WithRelease(*release)),
					1000)))

		var buf bytes.Buffer
		if err := validationEngine.Render(&buf, source, doc); err != nil {
			if len(documents) == 1 {
				fmt.Fprintf(clq.stderr, "❗️ %v\n", err)
			} else {
				fmt.Fprintf(clq.stderr, "❗️ %v: %v\n", document, err)
			}
			hasError = true
			continue
		}

		if buf.Len() > 0 {
			if len(documents) == 1 {
				fmt.Fprintln(clq.stdout, buf.String())
			} else {
				fmt.Fprintf(clq.stdout, "✅ %v: %v\n", document, buf.String())
			}
		}
	}

	if hasError {
		return 1
	}
	return 0
}

func (clq *Clq) readInput(input string) ([]byte, error) {
	if input == "-" {
		return ioutil.ReadAll(clq.stdin)
	}

	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
