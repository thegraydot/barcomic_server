MODULE := github.com/thegraydot/barcomic
BINARY := barcomic

.PHONY: fmt fmt_check mod_check vet test test_verbose test_coverage ci build install clean snapshot release_check

fmt:
	gofmt -w .

fmt_check:
	@test -z "$$(gofmt -l .)"

mod_check:
	go mod tidy && git diff --exit-code go.mod go.sum

vet:
	go vet ./...

test:
	go test -race -count=1 ./...

test_verbose:
	go test -race -count=1 -v ./...

test_coverage:
	go test -race -count=1 -coverpkg=./internal/... -coverprofile=coverage.out ./...

ci: fmt_check mod_check vet test

build:
	go build -ldflags="-s -w -X $(MODULE)/cmd.Version=dev" -o bin/$(BINARY) .

install:
	go install -ldflags="-s -w -X $(MODULE)/cmd.Version=dev" .

clean:
	rm -rf bin/ dist/

snapshot:
	goreleaser release --snapshot --clean

release_check:
	goreleaser check
