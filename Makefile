.PHONY: yacc
yacc:
	@go install golang.org/x/tools/cmd/goyacc@latest
	@cd internal/tags && goyacc tags.y

.PHONY: proto
proto:
	protoc -I include -I options --go_out=paths=source_relative:options options/*.proto
	protoc -I include -I options -I internal/replace/testdata --go_out=paths=source_relative:internal/replace/testdata internal/replace/testdata/msg.proto

.PHONY: test
test: 
	go test -race ./...

.PHONY: test-coverage
test-coverage: 
	go test -race -coverprofile=profile..txt -covermode=atomic ./...

.PHONY: build
build: 
	go build

.PHONY: install
install: 
	go install

