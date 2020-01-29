package tsm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func ReadUint64(reader io.Reader) (uint64, error) {

	var buffer = make([]byte, 8)
	_, err := reader.Read(buffer)
	if err != nil {
		return 0, fmt.Errorf("ReadUint64 Buffer Read: %w", err)
	}
	uintReader := bytes.NewReader(buffer)
	var data uint64
	err = binary.Read(uintReader, binary.BigEndian, &data)
	if err != nil {
		return 0, fmt.Errorf("ReadUint64 Binary Read: %w", err)
	}

	return data, nil
}

func ReadUint16(reader io.Reader) (uint16, error) {

	var buffer = make([]byte, 2)
	_, err := reader.Read(buffer)
	if err != nil {
		return 0, fmt.Errorf("ReadUint16 Buffer Read: %w", err)
	}
	uintReader := bytes.NewReader(buffer)
	var data uint16
	err = binary.Read(uintReader, binary.BigEndian, &data)
	if err != nil {
		return 0, fmt.Errorf("ReadUint16 Binary Read: %w", err)
	}

	return data, nil
}
