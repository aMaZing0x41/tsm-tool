package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/amazing0x41/tsm-tool/internal/tsm"
)

var fFile = flag.String("file", "", "tsm file to read")
var fDebug = flag.Bool("debug", false, "show additional debug logging")
var usageFunc = func() {
	fmt.Println(`./tsm-tool -file <file>`)
}

var readBytesFunc = func(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func main() {
	flag.Usage = usageFunc
	flag.Parse()

	if *fFile == "" {
		fmt.Println("Must provide a valid path to a tsm file.")
		usageFunc()
		return
	}

	absFilePath, err := filepath.Abs(*fFile)
	if err != nil {
		panic(err)
	}

	fileBytes, err := readBytesFunc(absFilePath)
	if err != nil {
		panic(err)
	}

	if *fDebug {
		fmt.Println("tsm file bytes:")
		fmt.Println(fileBytes)
	}

	data := readTSMFooter(&fileBytes)
	indexStart := int64(data.FooterPos)
	indexes := readTSMIndexes(indexStart, &fileBytes)
	crcErrors := []error{}
	fmt.Println("Reading TSM Blocks...")
	for i, index := range indexes {
		fmt.Println("***********************")
		fmt.Println("Block", i)
		fmt.Println("***********************")
		err := readTSMBlock(index.IndexEntry.BlockOffset, index.IndexEntry.BlockSize, &fileBytes)
		if err != nil {
			crcErrors = append(crcErrors, err)
		}
	}
	fmt.Println("Read TSM Blocks")

	if *fDebug {
		fmt.Println(data)
	}

	if len(crcErrors) > 0 {
		fmt.Println("CRC Validation FAILED")
		for _, e := range crcErrors {
			fmt.Println(e)
		}
		os.Exit(1)
	}
}

// readTSMFooter reades the footer information. This includes the start position for the indexes.
// returns a tsm.Footer with the positions populated. Panics if cannot read.
func readTSMFooter(fileBytes *[]byte) tsm.Footer {
	reader := bytes.NewReader(*fileBytes)
	data := tsm.Footer{}

	// TODO: these two lines should probably live somewhere else
	// this is easy b/c we have a reader and bytes
	fmt.Println("Reading TSM File")
	fmt.Println("File Size", reader.Size(), "bytes")

	// Read the footer to get the start position of the TSM index in the file
	pos, err := reader.Seek(-tsm.FooterSize, io.SeekEnd)
	if err != nil {
		fmt.Println("Read footer failed!")
		panic(err)
	}
	data.FooterIdx = pos
	_, err = reader.Seek(pos, io.SeekStart)
	if err != nil {
		fmt.Println("Failed to seek on footer")
		panic(err)
	}

	fpos, err := tsm.ReadUint64(reader)
	if err != nil {
		panic(err)
	}
	data.FooterPos = fpos
	return data
}

// TODO: convert all of the binary reads to a helper call in the tsm package (binary_reader)
// TODO: needs unit tests!
// readTSMIndexes seeks the reader to the start of the indexes, reads data, and moves the position to the end of that index.
// contiues until all indexs are read.
func readTSMIndexes(indexPos int64, fileBytes *[]byte) []tsm.IndexHeader {
	reader := bytes.NewReader(*fileBytes)

	// Seek to the begining of the indexes
	pos, err := reader.Seek(indexPos, io.SeekStart)
	if err != nil {
		panic(err)
	}

	indexes := []tsm.IndexHeader{}
	fmt.Println("Reading Indexes")
	for i := 1; pos < reader.Size()-tsm.FooterSize; i++ {
		fmt.Println("pos", pos)
		fmt.Println("***********************")
		fmt.Println("Index", i)
		fmt.Println("***********************")

		// Get key_len
		key_len, err := tsm.ReadUint16(reader)
		if err != nil {
			fmt.Println("Failed to get key_len")
			panic(err)
		}
		fmt.Println("key_len")
		fmt.Println(key_len)

		// Get key (string)
		var keyBuf = make([]byte, key_len)
		_, err = reader.Read(keyBuf)
		if err != nil {
			panic(err)
		}
		fmt.Println("key")
		fmt.Println(string(keyBuf))

		// Get type
		var idxTypeBuf = make([]byte, 1)
		_, err = reader.Read(idxTypeBuf)
		if err != nil {
			panic(err)
		}
		typeReader := bytes.NewReader(idxTypeBuf)
		var idxType byte
		err = binary.Read(typeReader, binary.BigEndian, &idxType)
		if err != nil {
			panic(err)
		}
		fmt.Println("Index Type")
		fmt.Println(idxType)

		// Get entry_count
		entry_count, err := tsm.ReadUint16(reader)
		if err != nil {
			fmt.Println("Failed to get entry_count")
			panic(err)
		}
		fmt.Println("entry_count")
		fmt.Println(entry_count)

		// Get min_time
		min_time, err := tsm.ReadUint64(reader)
		if err != nil {
			fmt.Println("Failed to read min_time")
			panic(err)
		}
		fmt.Println("min_time")
		fmt.Println(min_time)

		// Get max_time
		max_time, err := tsm.ReadUint64(reader)
		if err != nil {
			fmt.Println("Failed to read max_time")
			panic(err)
		}
		fmt.Println("max_time")
		fmt.Println(max_time)

		// Get block_offset
		bOffset, err := tsm.ReadUint64(reader)
		if err != nil {
			fmt.Println("Failed to read block offset")
			panic(err)
		}
		fmt.Println("block offset")
		fmt.Println(bOffset)

		// Get block_size
		var b_sizeBuf = make([]byte, 4)
		_, err = reader.Read(b_sizeBuf)
		if err != nil {
			panic(err)
		}
		bSizeReader := bytes.NewReader(b_sizeBuf)
		var bSize uint32
		err = binary.Read(bSizeReader, binary.BigEndian, &bSize)
		if err != nil {
			panic(err)
		}
		fmt.Println("Block Size")
		fmt.Println(bSize)

		pos, err = reader.Seek(0, io.SeekCurrent)
		if err != nil {
			fmt.Println("Failed to get index position")
			panic(err)
		}

		index := tsm.IndexHeader{
			KeyLen:     key_len,
			Key:        string(keyBuf),
			Type:       idxType,
			EntryCount: entry_count,
			IndexEntry: &tsm.IndexEntry{
				MinTime:     min_time,
				MaxTime:     max_time,
				BlockOffset: bOffset,
				BlockSize:   bSize,
			},
		}
		indexes = append(indexes, index)
	}
	fmt.Println("Read Indexes")
	return indexes
}

// readTSMBlock uses the information stored in the index to grab the stored
// CRC and compare that to the calculated CRC. If the CRC check fails, it returns
// an error; otherwise, nil
// TODO: needs unit tests!
func readTSMBlock(offset uint64, size uint32, fileBytes *[]byte) error {
	reader := bytes.NewReader(*fileBytes)
	_, err := reader.Seek(int64(offset), io.SeekStart)
	if err != nil {
		panic(err)
	}

	// Get CRC!!
	var crcBuf = make([]byte, 4)
	_, err = reader.Read(crcBuf)
	if err != nil {
		panic(err)
	}
	crcReader := bytes.NewReader(crcBuf)
	var crc uint32
	err = binary.Read(crcReader, binary.BigEndian, &crc)
	if err != nil {
		panic(err)
	}
	fmt.Println("stored CRC")
	fmt.Println(crc)

	var blockBuf = make([]byte, size-4)
	_, err = reader.Read(blockBuf)
	if err != nil {
		panic(err)
	}
	crcCalcd := crc32.ChecksumIEEE(blockBuf)
	fmt.Println("calculated CRC")
	fmt.Println(crcCalcd)

	if crc == crcCalcd {
		fmt.Println("Block Integrity Check PASSED")
	} else {
		fmt.Println("Block Integrity Check FAILED")
		return fmt.Errorf("Block Integrity Check FAILED - Stored CRC %v != Calculated CRC %v", crc, crcCalcd)
	}

	return nil
}
