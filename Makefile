GO=go
BINARY=./mdexec

all: $(BINARY) README.md

$(BINARY): cmd/mdexec/main.go
	$(GO) build -o $@ $<

README.md: src/README.md
	$(BINARY) $< > $@

clean:
	rm -f $(BINARY)
	rm -f README.md