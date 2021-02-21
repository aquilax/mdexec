package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/acarl005/stripansi"
	"github.com/aquilax/mdexec"
)

const codefence = "```"

const defaultTemplate = codefence + `sh
$ {{ .Command }}
{{ stripAnsi .Output }}
` + codefence + `
`

func getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"stripAnsi": func(output string) string {
			return stripansi.Strip(output)
		},
	}
}

func main() {
	blockTemplate := flag.String("template", defaultTemplate, "Template to use when rendering a command block")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] [FILE]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Execute commands in markdown and embeds the result in the output")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "FILE can be both file name or - to read from stdin")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "Fields available in the template:")
		fmt.Fprintln(flag.CommandLine.Output(), "  .Command  string - The command that was executed")
		fmt.Fprintln(flag.CommandLine.Output(), "  .Output   string - Command output")
		fmt.Fprintln(flag.CommandLine.Output(), "  .Error    error  - Execution error")
		fmt.Fprintln(flag.CommandLine.Output(), "  .Duration int64  - Execution duration in ns")
		fmt.Fprintln(flag.CommandLine.Output(), "Template functions:")
		fmt.Fprintln(flag.CommandLine.Output(), "  stripAnsi var    - Strips the ansi characters from the variable")
	}

	flag.Parse()
	var stream io.ReadCloser
	var err error
	if len(os.Args) > 1 && os.Args[1] != "" {
		if os.Args[1] == "-" {
			stream = os.Stdin
		} else {
			stream, err = os.Open(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			defer stream.Close()

		}
		template := template.Must(template.New("template").Funcs(getTemplateFunctions()).Parse(*blockTemplate))
		err = mdexec.ProcessStream(stream, os.Stdout, template)
		if err != nil {
			log.Fatal(err)
		}
	}
}
