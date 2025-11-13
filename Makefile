.PHONY: build run clean dev templ install

build: templ
	@go build -o openping ./cmd/openping

run: templ
	@go run ./cmd/openping

dev:
	@air & bun watch

templ:
	@templ generate

clean:
	@rm -f openping
	@rm -rf tmp/

install: build
	@mv openping /usr/local/bin/
