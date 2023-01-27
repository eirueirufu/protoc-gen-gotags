package options

//go:generate protoc -I . -I ../include --go_out=. --go_opt=paths=source_relative ./*.proto
