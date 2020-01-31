package tsm

// TODO: Need to understand how IndexEntries are constructed.
// Can there be more than one of them in an index header?
// When is that true?
// Doc is unclear
type IndexHeader struct {
	KeyLen     uint16
	Key        string
	Type       byte
	EntryCount uint16
	IndexEntry *IndexEntry
}

type IndexEntry struct {
	MinTime     uint64
	MaxTime     uint64
	BlockOffset uint64
	BlockSize   uint32
}

type BlockEntry struct {
	CRC32 uint32
	Data  []byte
}
