package replace

import (
	"crypto/md5"
	"os"
	"reflect"
	"testing"
)

func TestReplacer_replaceTags(t *testing.T) {
	src, err := os.ReadFile("./testdata/msg.pb.go")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("./testdata/want")
	if err != nil {
		t.Fatal(err)
	}
	replace := NewReplacer(".")
	replace.msg = map[string]msg{
		"Msg1": {
			"NoTagOption": {},
			"AppendTag": {
				part: map[string]string{
					"append_key": "msg1_append_value",
				},
			},
			"ReplaceTag": {
				part: map[string]string{
					"json": "msg1_replace_json_value",
				},
			},
			"AppendTags": {
				part: map[string]string{
					"append_key1": "msg1_append_value1",
					"append_key2": "msg1_append_value2",
				},
			},
			"ReplaceTags": {
				part: map[string]string{
					"protobuf": "msg1_replace_proto_value",
					"json":     "msg1_replace_json_value",
				},
			},
			"ReplaceAndAppendTags": {
				part: map[string]string{
					"json":       "msg1_replace_json_value",
					"append_key": "msg1_append_value",
				},
			},
			"All": {
				all: "all:\"msg1_all\"",
			},
			"AllWithReplaceAndAppendTags": {
				all: "all:\"msg1_all\"",
				part: map[string]string{
					"json":       "msg1_replace_json_value",
					"append_key": "msg1_append_value",
				},
			},
		},
		"Msg2": {
			"NoTagOption": {},
			"AppendTag": {
				part: map[string]string{
					"append_key": "msg2_append_value",
				},
			},
			"ReplaceTag": {
				part: map[string]string{
					"json": "msg2_replace_json_value",
				},
			},
			"AppendTags": {
				part: map[string]string{
					"append_key1": "msg2_append_value1",
					"append_key2": "msg2_append_value2",
				},
			},
			"ReplaceTags": {
				part: map[string]string{
					"protobuf": "msg2_replace_proto_value",
					"json":     "msg2_replace_json_value",
				},
			},
			"ReplaceAndAppendTags": {
				part: map[string]string{
					"json":       "msg2_replace_json_value",
					"append_key": "msg2_append_value",
				},
			},
			"All": {
				all: "all:\"msg2_all\"",
			},
			"AllWithReplaceAndAppendTags": {
				all: "all:\"msg2_all\"",
				part: map[string]string{
					"json":       "msg2_replace_json_value",
					"append_key": "msg2_append_value",
				},
			},
		},
	}

	got, err := replace.replaceTags("msg.pb.go", src)
	if err != nil {
		t.Fatal(err)
	}

	gotMd5, wantMd5 := md5.Sum(got), md5.Sum(want)
	if !reflect.DeepEqual(gotMd5, wantMd5) {
		t.Fatalf("Replacer.replaceTags() = %x, want %x", gotMd5, wantMd5)
	}
}
