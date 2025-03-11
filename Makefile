# Zachary Perry
# compiles problem_1 and problem_2 into respective bin/ files
.DEFAULT_GOAL := build

fmt: 
	@go fmt ./...

lint: fmt
	@golint ./...

vet: fmt
	@go vet ./...

build: vet
	@go build -o bin/problem_1 cmd/problem_1/problem_1.go
	@go build -o bin/problem_2 cmd/problem_2/problem_2.go

clean:
	@go clean
	rm bin/*

.PHONY: fmt lint vet build clean
