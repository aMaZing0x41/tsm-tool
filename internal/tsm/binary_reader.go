package tsm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func ReadUint64AtPos(reader io.ReaderAt, off int64) (uint64, error) {

	var buffer = make([]byte, 8)
	_, err := reader.ReadAt(buffer, off)
	if err != nil {
		return 0, fmt.Errorf("ReadUint64AtPos ReadAt: %w", err)
	}
	uintReader := bytes.NewReader(buffer)
	var data uint64
	err = binary.Read(uintReader, binary.BigEndian, &data)
	if err != nil {
		return 0, fmt.Errorf("ReadUint64AtPos Read: %w", err)
	}

	fmt.Println(data)
	return data, nil
}
