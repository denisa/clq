package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/query"
	"github.com/denisa/clq/internal/validator"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var version string

type Clq struct {
	stdin          io.Reader
	stdout, stderr io.Writer
	verbose        bool
	documents      []string
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
	var changeMap = options.String("changeMap", "", "Name of a file defining the mapping from change kind to semantic version change")
	var output = options.String("output", "json", "Output format, for complex result. One of: json|md")
	var queryString = options.String("query", "", "A query to extract information out of the change log")
	var release = options.Bool("release", false, "Enable release-mode validation")
	var showVersion = options.Bool("version", false, "Prints clq version")
	options.BoolVar(&clq.verbose, "with-filename", false, "Always print filename headers with output lines")

	if options.Parse(arguments) != nil {
		return 2
	}

	if *showVersion {
		fmt.Fprintf(clq.stdout, "clq %v\n", version)
		return 0
	}

	if options.NArg() == 0 {
		clq.documents = []string{"-"}
	} else {
		clq.documents = options.Args()
	}

	var hasError bool
	for _, document := range clq.documents {
		queryEngine, err := query.NewQueryEngine(*queryString, *output)
		if err != nil {
			clq.error("", err)
			return 2
		}
		changeKind, err := changelog.NewChangeKind(*changeMap)
		if err != nil {
			clq.error("", err)
			return 2
		}
		source, err := clq.readInput(document)
		if err != nil {
			clq.error(document, err)
			hasError = true
			continue
		}

		reader := text.NewReader(source)
		doc := goldmark.DefaultParser().Parse(reader)

		validatorOpts := []validator.Option{validator.WithRelease(*release), validator.WithChangeKind(changeKind)}
		if queryEngine.HasQuery() {
			validatorOpts = append(validatorOpts, validator.WithListener(queryEngine))
		}
		validationEngine := renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(validator.NewValidator(validatorOpts...), 1000)))

		var buf bytes.Buffer
		if err := validationEngine.Render(&buf, source, doc); err != nil {
			clq.error(document, err)
			hasError = true
			continue
		}
		clq.output(document, queryEngine.Result())
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
	return ioutil.ReadFile(input)
}

func (clq *Clq) withFileName() bool {
	return clq.verbose || len(clq.documents) > 1
}

func (clq *Clq) error(document string, err error) {
	if err, ok := err.(*os.PathError); ok {
		fmt.Fprintf(clq.stderr, "❗️ %v: %v\n", err.Path, err.Err.Error())
		return
	}

	if clq.withFileName() && document != "" {
		fmt.Fprintf(clq.stderr, "❗️ %v: %v\n", document, err)
	} else {
		fmt.Fprintf(clq.stderr, "❗️ %v\n", err)
	}
}

func (clq *Clq) output(document string, result string) {
	if result == "" {
		if clq.verbose {
			fmt.Fprintf(clq.stdout, "✅ %v\n", document)
		}
	} else if clq.withFileName() {
		fmt.Fprintf(clq.stdout, "✅ %v: %v\n", document, result)
	} else {
		fmt.Fprintln(clq.stdout, result)
	}
}
