package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	texttmpl "text/template"

	templatesvc "hufschlaeger.net/markscribe/internal/service/template"
)

var (
	write = flag.String("write", "", "write output to")
)

func main() {
	flag.Parse()

	// Support placing -write after the template argument by scanning remaining args.
	args := flag.Args()
	var (
		templatePath  string
		writeOverride string
	)
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-write" || a == "--write":
			if i+1 >= len(args) {
				fmt.Println("Missing value for -write")
				os.Exit(1)
			}
			writeOverride = args[i+1]
			i++ // consume value
		case strings.HasPrefix(a, "-write="):
			writeOverride = strings.TrimPrefix(a, "-write=")
		case strings.HasPrefix(a, "--write="):
			writeOverride = strings.TrimPrefix(a, "--write=")
		default:
			if templatePath == "" {
				templatePath = a
			}
		}
	}

	if templatePath == "" {
		fmt.Println("Usage: markscribe [template] [-write output]\nExamples:\n  markscribe README.md.tpl\n  markscribe README.md.tpl -write README.md")
		os.Exit(1)
	}

	if writeOverride != "" && *write == "" {
		// allow override only if not already set via flags
		*write = writeOverride
	}

	tplIn, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Can't read file:", err)
		os.Exit(1)
	}

	// Build template service from environment to keep startup lean
	tplSvc, err := templatesvc.NewFromEnv(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create template with template service functions to declutter main
	tpl, err := texttmpl.New("tpl").Funcs(tplSvc.Funcs()).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
	}

	w := os.Stdout

	if len(*write) > 0 {
		f, err := os.Create(*write)
		if err != nil {
			fmt.Println("Can't create:", err)
			os.Exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println("Can't close file:", err)
			}
		}(f)
		w = f
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Can't render template:", err)
		os.Exit(1)
	}
}
