package testdata

//go:generate protoc -I . -I ../../../include -I ../../../options --go_out=. --go_opt=paths=source_relative ./*.proto
