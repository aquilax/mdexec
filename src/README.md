# mdexec

[![Go Reference](https://pkg.go.dev/badge/github.com/aquilax/mdexec.svg)](https://pkg.go.dev/github.com/aquilax/mdexec)

Executes commands in a markdown file and embeds the result

## Format

Prefix the commands as inline blocks e.g.:

```markdown
`#$ cat examples/log.yaml`
```

Which will embed the content of examples/log.yaml in the output markdown text.

If you want to prevent the execution use `#$` instead of `$`.
The hash `#` symbol will be removed from the output

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

## Usage

`$ ./mdexec -help`

## Similar projects

* https://zimbatm.github.io/mdsh/
