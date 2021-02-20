# mdexec
Executes commands in a markdown file and embeds the result

## Usage

Run the command from the working directory

```
$ mdexec md_template_file.md > output.md
```

or

```
$ mdexec < md_template_file.md > output.md
```

Prefix the commands as inline blocks e.g.:

```
`$ cat examples/log.yaml`
```

Which will embed the content of examples/log.yaml in the output markdown text.