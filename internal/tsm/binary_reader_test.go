package tsm

import (
	"bytes"
	"io"
	"math"
	"testing"
)

func getBytesReader(b []byte, off int64) *bytes.Reader {
	r := bytes.NewReader(b)
	r.Seek(off, io.SeekStart)
	return r
}

func TestReadUint64(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			"zero",
			args{
				reader: getBytesReader([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, 0),
			},
			uint64(0),
			false,
		},
		{
			"max",
			args{
				reader: getBytesReader([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0),
			},
			math.MaxUint64,
			false,
		},
		{
			"offset",
			args{
				reader: getBytesReader([]byte{0xf, 0xf, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11}, 2),
			},
			uint64(17),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadUint64(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint64AtPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint64AtPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadUint16(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			"zero",
			args{
				reader: getBytesReader([]byte{0x0, 0x0}, 0),
			},
			uint16(0),
			false,
		},
		{
			"max",
			args{
				reader: getBytesReader([]byte{0xff, 0xff}, 0),
			},
			math.MaxUint16,
			false,
		},
		{
			"offset",
			args{
				reader: getBytesReader([]byte{0xf, 0xf, 0x0, 0x11}, 2),
			},
			uint16(17),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadUint16(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}
