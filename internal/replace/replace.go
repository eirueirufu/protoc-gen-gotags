package replace

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	"github.com/eirueirufu/protoc-gen-gotags/internal/tags"
	"github.com/eirueirufu/protoc-gen-gotags/options"
)

type (
	Replacer struct {
		msg     map[string]msg
		fileSet *token.FileSet
	}
	msg map[string]tag
	tag struct {
		all  string
		part map[string]string
	}

	tagKv struct {
		key, val string
	}

	placeholderType struct{}
)

func NewReplacer() *Replacer {
	return &Replacer{
		msg:     map[string]msg{},
		fileSet: token.NewFileSet(),
	}
}

func newMsg() msg {
	return msg{}
}

func newTag() tag {
	return tag{
		part: map[string]string{},
	}
}

func (p *Replacer) replaceTags(filename string, src []byte) ([]byte, error) {
	astFile, err := parser.ParseFile(p.fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			msgName := typeSpec.Name.Name
			structType := typeSpec.Type.(*ast.StructType)
			for _, field := range structType.Fields.List {
				if len(field.Names) != 1 {
					continue
				}
				fieldName := field.Names[0].Name
				tagVal := ""
				if field.Tag != nil {
					tagVal = field.Tag.Value
				}
				replacedTag, err := p.replaceTag(msgName, fieldName, tagVal)
				if err != nil {
					return nil, err
				}
				if len(replacedTag) == 0 {
					continue
				}
				if field.Tag == nil {
					field.Tag = &ast.BasicLit{
						Kind: token.STRING,
					}
				}
				field.Tag.Value = replacedTag
			}
		}
	}
	buff := &bytes.Buffer{}
	if err := format.Node(buff, p.fileSet, astFile); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (p *Replacer) replaceTag(msgName, fieldName, tagVal string) (string, error) {
	msg, ok := p.msg[msgName]
	if !ok {
		return tagVal, nil
	}
	fieldTag, ok := msg[fieldName]
	if !ok {
		return tagVal, nil
	}
	if len(fieldTag.all) > 0 {
		return fmt.Sprintf("`%s`", fieldTag.all), nil
	}

	tags, err := tags.ParseTags(tagVal)
	if err != nil {
		return "", fmt.Errorf("tags error on message:%s field:%s, details:%s", msgName, fieldName, err.Error())
	}

	tagKvs := make([]tagKv, 0)
	marked := map[string]placeholderType{}
	for _, tag := range tags {
		key := tag.Key
		val, ok := fieldTag.part[key]
		if !ok {
			val = tag.Value
		}
		tagKvs = append(tagKvs, tagKv{
			key: key,
			val: val,
		})
		marked[key] = placeholderType{}
	}
	appendTagKvs := make([]tagKv, 0)
	for key, val := range fieldTag.part {
		if _, ok := marked[key]; ok {
			continue
		}
		appendTagKvs = append(appendTagKvs, tagKv{
			key: key,
			val: val,
		})
	}
	sort.Slice(appendTagKvs, func(i, j int) bool {
		return appendTagKvs[i].key < appendTagKvs[j].key
	})
	tagKvs = append(tagKvs, appendTagKvs...)
	tagStrs := make([]string, 0)
	for _, tagKv := range tagKvs {
		tagStrs = append(tagStrs, fmt.Sprintf("%s:\"%s\"", tagKv.key, tagKv.val))
	}
	return fmt.Sprintf("`%s`", strings.Join(tagStrs, " ")), nil
}

func (p *Replacer) ParseFile(file *protogen.File) error {
	filename := file.GeneratedFilenamePrefix + ".pb.go"
	src, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("make sure you have used the path type 'source_relative', %w", err)
		}
		return err
	}
	for _, msg := range file.Messages {
		if err := p.parseMsg(msg); err != nil {
			return err
		}
	}
	bs, err := p.replaceTags(filename, src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, bs, 0755); err != nil {
		return err
	}
	return nil
}

func (p *Replacer) parseMsg(msg *protogen.Message) error {
	msgName := msg.GoIdent.GoName
	p.msg[msgName] = newMsg()
	for _, field := range msg.Fields {
		if err := p.parseField(field, p.msg[msgName]); err != nil {
			return err
		}
	}
	return nil
}

func (p *Replacer) parseField(field *protogen.Field, mp msg) error {
	fieldName := field.GoName
	tag := newTag()

	opts := field.Desc.Options()
	if proto.HasExtension(opts, options.E_Tag) {
		ext := proto.GetExtension(opts, options.E_Tag)
		gotags, ok := ext.([]*options.Tag)
		if !ok {
			return fmt.Errorf("extension is %T; want a []Gotag", ext)
		}
		for _, gotag := range gotags {
			tag.part[gotag.GetKey()] = gotag.GetValue()
		}
	}
	if proto.HasExtension(opts, options.E_All) {
		ext := proto.GetExtension(opts, options.E_All)
		all, ok := ext.(string)
		if !ok {
			return fmt.Errorf("extension is %T; want a string", ext)
		}
		tag.all = all
	}
	mp[fieldName] = tag
	return nil
}
