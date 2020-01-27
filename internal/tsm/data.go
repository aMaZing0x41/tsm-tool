package tsm

import "fmt"

type TSM struct {
	FooterIdx int64
	FooterPos uint64
}

func (t TSM) String() string {
	return fmt.Sprintf("%v", t.FooterIdx)
}
