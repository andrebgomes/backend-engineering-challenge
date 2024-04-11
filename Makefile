.PHONY: buil test coverage coverage-html bench

build:
	go build -o unbabel_cli cmd/unbabel_cli/main.go 

test:
	go test -v -count=1 ./...

coverage:
	go test -coverprofile cover.out ./...

coverage-html: coverage
	go tool cover -html cover.out -o cover.html
	open cover.html

bench:
	go test -bench=. ./internal/translationdeliverytime