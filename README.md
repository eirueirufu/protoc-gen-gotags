# protoc-gen-gotags

<p align="left">
<a href="https://github.com/eirueirufu/protoc-gen-gotags/actions"><img src="https://github.com/eirueirufu/protoc-gen-gotags/workflows/go/badge.svg?branch=main" alt="Build Status"></a>
<a href="https://codecov.io/github/eirueirufu/protoc-gen-gotags"><img src="https://codecov.io/github/eirueirufu/protoc-gen-gotags/branch/main/graph/badge.svg?token=5NW5CP5H6G" alt="codeCov"></a>
</p>

## What is this?

In [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go), it is difficult to customize tags such as gorm, db, etc. You can use this tool to add any tag you need to the go field.

## Install

Download in [release page](https://github.com/eirueirufu/protoc-gen-gotags/releases) or use `go install`

```golang
go install github.com/eirueirufu/protoc-gen-gotags@latest
```

## Usage

1. define the options
2. use `protoc-gen-go` to generate your `.pb.go` file
3. use `protoc-gen-gotags` to replace the generated `.pb.go` file tags, in this cmd, you should set opt param `go_out` to specify protoc-gen-go out_dir, such as `--gotags_opt=go_out=.` 

field options example：

```proto
syntax = "proto3";

package example;

import "options.proto";

......

message Msg {
    string append_tag = 1 [(gotags.tag) = {
        key: "append_key",
        value: "msg_append_value",
    }];
    string replace_tag = 2 [(gotags.tag) = {
        key: "json",
        value: "msg_replace_json_value",
    }];
    string append_tags = 3 [
        (gotags.tag) = {
            key: "append_key1",
            value: "msg_append_value1",
        },
        (gotags.tag) = {
            key: "append_key2",
            value: "msg_append_value2",
        }
    ];
    string replace_tags = 4 [
        (gotags.tag) = {
            key: "protobuf",
            value: "msg_replace_proto_value",
        },
        (gotags.tag) = {
            key: "json",
            value: "msg_replace_json_value",
        }
    ];
    string replace_and_append_tags = 5 [
        (gotags.tag) = {
            key: "json",
            value: "msg_replace_json_value",
        },
        (gotags.tag) = {
            key: "append_key",
            value: "msg_append_value",
        }
    ];
    string all = 6 [(gotags.all) = "all:\"msg_all\""];
}
```

before:

```go
type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppendTag            string `protobuf:"bytes,1,opt,name=append_tag,json=appendTag,proto3" json:"append_tag,omitempty"`
	ReplaceTag           string `protobuf:"bytes,2,opt,name=replace_tag,json=replaceTag,proto3" json:"replace_tag,omitempty"`
	AppendTags           string `protobuf:"bytes,3,opt,name=append_tags,json=appendTags,proto3" json:"append_tags,omitempty"`
	ReplaceTags          string `protobuf:"bytes,4,opt,name=replace_tags,json=replaceTags,proto3" json:"replace_tags,omitempty"`
	ReplaceAndAppendTags string `protobuf:"bytes,5,opt,name=replace_and_append_tags,json=replaceAndAppendTags,proto3" json:"replace_and_append_tags,omitempty"`
	All                  string `protobuf:"bytes,6,opt,name=all,proto3" json:"all,omitempty"`
}
```

after:

```go
type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppendTag            string `protobuf:"bytes,1,opt,name=append_tag,json=appendTag,proto3" json:"append_tag,omitempty" append_key:"msg_append_value"`
	ReplaceTag           string `protobuf:"bytes,2,opt,name=replace_tag,json=replaceTag,proto3" json:"msg_replace_json_value"`
	AppendTags           string `protobuf:"bytes,3,opt,name=append_tags,json=appendTags,proto3" json:"append_tags,omitempty" append_key1:"msg_append_value1" append_key2:"msg_append_value2"`
	ReplaceTags          string `protobuf:"msg_replace_proto_value" json:"msg_replace_json_value"`
	ReplaceAndAppendTags string `protobuf:"bytes,5,opt,name=replace_and_append_tags,json=replaceAndAppendTags,proto3" json:"msg_replace_json_value" append_key:"msg_append_value"`
	All                  string `all:"msg_all"`
}
```

## License

See [LICENSE](./LICENSE).