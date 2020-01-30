package main

import "testing"

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
