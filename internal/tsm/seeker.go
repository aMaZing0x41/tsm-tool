package tsm

import (
	"fmt"
	"io"
)

// SeekFromEnd take an io.Seeker and offset and seeks the reader from the end to the offset.
// The position is returned along with an error.
func SeekFromEnd(reader io.Seeker, offset int64) (int64, error) {
	pos, err := reader.Seek(offset, io.SeekEnd)
	if err != nil {
		return 0, fmt.Errorf("could not seek: %w", err)
	}

	return pos, nil
}
