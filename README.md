# mdexec

Executes commands in a markdown file and embeds the result

## Format

Prefix the commands as inline blocks e.g.:

```markdown
`$ cat examples/log.yaml`
```

Which will embed the content of examples/log.yaml in the output markdown text.

## Installation

```sh
go install github.com/aquilax/mdexec/cmd/mdexec
```

## Cmd Usage

Run the command from the working directory

```sh
mdexec md_template_file.md > output.md
```

or

```sh
cat md_template_file.md | mdexec > output.md
```


```sh
mdexec -h
Usage: mdexec [OPTIONS] [FILE]
Execute commands in markdown and embeds the result in the output

FILE can be both file name or - to read from stdin

  -template string
        Template to use when rendering a command block (default "```sh\n$ {{ .Command }}\n{{ .Output }}\n```\n")

Fields available in the template:
  {{ .Command }}  string - the command that was executed
  {{ .Output }}   string - command output
  {{ .Error }}    error  - Execution error
  {{ .Duration }} int64  - execution duration in ns
```