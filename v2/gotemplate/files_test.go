package gotemplate

import (
	"reflect"
	"testing"
)

func TestFindFiles(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "FindFiles",
			args: args{path: "test_files"},
			want: []string{
				"test_files/10/20/file_21",
				"test_files/file_01",
				"test_files/file_02",
			},
			wantErr: false,
		},
		{
			name:    "Invalid Folder",
			args:    args{path: "invlaid/folder"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindFiles(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
