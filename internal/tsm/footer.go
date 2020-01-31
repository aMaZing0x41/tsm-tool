package tsm

import "fmt"

type Footer struct {
	FooterIdx int64
	FooterPos uint64
}

func (t Footer) String() string {
	return fmt.Sprintf("TSM Footer Information\n\tFooterIdx: %v,\n\tFooterPos: %v", t.FooterIdx, t.FooterPos)
}
