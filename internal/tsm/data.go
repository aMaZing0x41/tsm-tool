package tsm

import "fmt"

type TSM struct {
	FooterIdx int64
	FooterPos uint64
}

func (t TSM) String() string {
	return fmt.Sprintf("TSM Information\n\tFooterIdx: %v,\n\tFooterPos: %v", t.FooterIdx, t.FooterPos)
}
