.PHONY: buil test coverage coverage-html

build:
	go build -o unbabel_cli cmd/unbabel_cli/main.go 

test:
	go test -short -count=1 ./... || exit 1

coverage:
	go test -v -coverprofile cover.out ./...

coverage-html:
	go tool cover -html cover.out -o cover.html
	open cover.html
