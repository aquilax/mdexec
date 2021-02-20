package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
)

func execute(command string) (string, error) {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	args, err := shellquote.Split(command)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = wd
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func processFile(file io.Reader) {
	scanner := bufio.NewScanner(file)
	var line string
	var trimmedLine string
	var command string
	var output string
	var err error

	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "`$ ") {
			trimmedLine = strings.TrimSpace(line)
			// get command
			command = trimmedLine[3 : len(trimmedLine)-1]
			// run command
			output, err = execute(command)
			if err != nil {
				log.Fatalln("Error running command: `" + command + "` " + err.Error())
			}
			// generate output
			fmt.Fprintln(os.Stdout, "```sh")
			fmt.Fprintln(os.Stdout, "$ "+command)
			fmt.Fprintln(os.Stdout, output)
			fmt.Fprintln(os.Stdout, "```")
			continue
		}
		fmt.Fprintln(os.Stdout, line)
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		processFile(file)
	} else {
		processFile(os.Stdin)
	}
}
