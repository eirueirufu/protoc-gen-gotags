package main

import (
	"flag"
	"fmt"

	"github.com/eirueirufu/protoc-gen-gotags/internal/replace"
	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.1.2"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-gotag %v\n", version)
		return
	}

	var flags flag.FlagSet
	goOut := flags.String("go_out", ".", "must set the same value as go_out, such as '--gotags_opt=go_out=.'")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		replacer := replace.NewReplacer(*goOut)
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
