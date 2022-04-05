package mdexec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/kballard/go-shellquote"
)

const codefence = "```"

const defaultTemplate = codefence + `sh
$ {{ .Command }}
{{ .Output }}
` + codefence + `
`

type commandContext struct {
	workDir string
	env     []string
}

// TemplateContext contains all fields available in the template
type TemplateContext struct {
	Command  string
	Output   string
	Duration int64
	Error    error
}

// Executor is a function that when given command must return the output, the duration and optional error
type Executor func(command string) (string, int64, error)

func getDefaultExecutor() (Executor, error) {
	var err error
	var context commandContext
	context.env = os.Environ()
	context.workDir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	return func(command string) (string, int64, error) {
		var err error
		args, err := shellquote.Split(command)
		if err != nil {
			return "", 0, err
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = context.env
		cmd.Dir = context.workDir
		var out bytes.Buffer
		cmd.Stdout = &out
		start := time.Now()
		err = cmd.Run()
		duration := time.Since(start)
		if err != nil {
			return "", duration.Nanoseconds(), err
		}
		return out.String(), duration.Nanoseconds(), nil
	}, nil
}

// ProcessStream processes commands in the inStream rendering the tmpl template and writing the output to the outStream
func ProcessStream(inStream io.Reader, outStream io.Writer, tmpl *template.Template) error {
	executor, err := getDefaultExecutor()
	if err != nil {
		return err
	}
	return ProcessStreamWithExecutor(inStream, outStream, tmpl, executor)
}

// ProcessStreamWithExecutor processes commands in the inStream rendering the tmpl template and writing the output to the outStream
// calling the executor function for each command
func ProcessStreamWithExecutor(inStream io.Reader, outStream io.Writer, tmpl *template.Template, executor Executor) error {
	scanner := bufio.NewScanner(inStream)
	var line string
	var trimmedLine string
	var err error
	var templateContext TemplateContext
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "`$ ") {
			trimmedLine = strings.TrimSpace(line)
			// get command
			templateContext.Command = trimmedLine[3 : len(trimmedLine)-1]
			// run command
			templateContext.Output, templateContext.Duration, err = executor(templateContext.Command)
			if err != nil {
				templateContext.Error = err
				log.Println("Error running command: `" + templateContext.Command + "` " + err.Error())
			}
			// generate output
			tmpl.Execute(outStream, templateContext)
			continue
		}
		if strings.HasPrefix(line, "`#$ ") {
			line = string(append([]byte{'`', '$'}, []byte(line)[4:]...))
		}
		fmt.Fprintln(outStream, line)
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}
