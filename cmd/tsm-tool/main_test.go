package main

import (
	"reflect"
	"testing"

	"github.com/amazing0x41/tsm-tool/internal/tsm"
)

func Test_main(t *testing.T) {
	readBytesFunc = func(x string) ([]byte, error) {
		return testFileBytes, nil
	}
	tests := []struct {
		name      string
		bytesFunc func(string) ([]byte, error)
	}{
		{"default", readBytesFunc},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_readTSMFooter(t *testing.T) {
	type args struct {
		fileBytes *[]byte
	}
	tests := []struct {
		name string
		args args
		want tsm.Footer
	}{
		{"test_1", args{fileBytes: &testFileBytes}, tsm.Footer{FooterPos: 4696, FooterIdx: 4882}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readTSMFooter(tt.args.fileBytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTSMFooter() = %v, want %v", got, tt.want)
			}
		})
	}
}
