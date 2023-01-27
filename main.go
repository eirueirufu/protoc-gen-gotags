package main

import (
	"flag"
	"fmt"

	"github.com/eirueirufu/protoc-gen-gotags/internal/replace"
	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-gotag %v\n", version)
		return
	}

	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		replacer := replace.NewReplacer()
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if err := replacer.ParseFile(f); err != nil {
				return err
			}
		}
		return nil
	})
}
