install-air:
	go install github.com/air-verse/air@latest

run-debug @install-air:
	air --build.cmd "go build -o ./bin/http ./cmd/http/main.go" --build.bin "./bin/http"

build:
	go build -o ./bin/http ./cmd/http/main.go

run: build
	bin/http