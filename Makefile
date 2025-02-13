install-air:
	go install github.com/air-verse/air@latest

run-debug @install-air:
	air --build.cmd "go build -o ./bin/http ./cmd/http/main.go" --build.bin "./bin/http"

run-consumers-debug @install-air:
	air --build.cmd "go build -o ./bin/consumers ./cmd/consumers/main.go" --build.bin "./bin/consumers"

build:
	go build -o ./bin/http ./cmd/http/main.go

run: build
	bin/http