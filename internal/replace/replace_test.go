package replace

import (
	"reflect"
	"testing"
)

func TestReplacer_replaceTags(t *testing.T) {
	type args struct {
		filename string
		src      []byte
	}
	tests := []struct {
		name    string
		p       *Replacer
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.replaceTags(tt.args.filename, tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replacer.replaceTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replacer.replaceTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
