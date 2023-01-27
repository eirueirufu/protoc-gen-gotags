package tags

import (
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    Tags
		wantErr bool
	}{
		{
			name: "no tags01",
			args: args{
				src: "",
			},
		},
		{
			name: "no tags02",
			args: args{
				src: "``",
			},
		},
		{
			name: "no key",
			args: args{
				src: "`:\"foo\"`",
			},
			wantErr: true,
		},
		{
			name: "no val",
			args: args{
				src: "`foo:\"\"`",
			},
			want: Tags{{Key: "foo", Value: ""}},
		},
		{
			name: "no colon",
			args: args{
				src: "`foo\"\"`",
			},
			wantErr: true,
		},
		{
			name: "single quote",
			args: args{
				src: "`foo\"`",
			},
			wantErr: true,
		},
		{
			name: "tag",
			args: args{
				src: "`tag:\"foo\"`",
			},
			want: Tags{{Key: "tag", Value: "foo"}},
		},
		{
			name: "tag with whitespace",
			args: args{
				src: "`\t\n\r tag\t\n\r :\t\n\r \"foo\"\t\n\r `",
			},
			want: Tags{{Key: "tag", Value: "foo"}},
		},
		{
			name: "tags",
			args: args{
				src: "`foo1:\"bar1\" foo2:\"bar2\"`",
			},
			want: Tags{{Key: "foo1", Value: "bar1"}, {Key: "foo2", Value: "bar2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTags(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParseTags(b *testing.B) {
	src := "`tag1:\"foo\" tag2:\"bar\"`"
	for i := 0; i < b.N; i++ {
		_, err := ParseTags(src)
		if err != nil {
			b.Fatal(err)
		}
	}
}
