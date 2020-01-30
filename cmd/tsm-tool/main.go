package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	readTSMIndexes(indexStart, &fileBytes)

	if *fDebug {
		fmt.Println(data)
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
	reader.Seek(pos, io.SeekStart)

	fpos, err := tsm.ReadUint64(reader)
	if err != nil {
		panic(err)
	}
	data.FooterPos = fpos
	return data
}

// TODO: convert all of the binary reads to a helper call in the tsm package (binary_reader)
// TODO: create an index struct and store the information for each index so that we can use it later to iterate over the blocks
// TODO: needs unit tests!
// readTSMIndexes seeks the reader to the start of the indexes, reads data, and moves the position to the end of that index.
// contiues until all indexs are read.
func readTSMIndexes(indexPos int64, fileBytes *[]byte) {
	reader := bytes.NewReader(*fileBytes)

	// Seek to the begining of the indexes
	pos, err := reader.Seek(indexPos, io.SeekStart)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading TSM Indexes")
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

		// Get key
		var keyBuf = make([]byte, key_len)
		_, err = reader.Read(keyBuf)
		if err != nil {
			panic(err)
		}
		fmt.Println("key")
		fmt.Println(string(keyBuf))

		// Advance over type, entry count, min and max time
		offset := int64(0)
		offset += 1 // type
		offset += 2 // entry count
		offset += 8 // min_time
		offset += 8 // max_time
		offset, err = reader.Seek(offset, io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		fmt.Println("position")
		fmt.Println(offset)

		// Get block_offset
		bOffset, err := tsm.ReadUint64(reader)
		if err != nil {
			fmt.Println("Failed to read block offset.")
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
	}
	fmt.Println("Read TSM Indexes")
}

// TODO: before this func can be lit up, the readTSMIndexes func must collect and store all of the information
// about each index. We can then iterate over the index information and read each block and get the CRC etc.
func readTSMBlock(offset int64, fileBytes *[]byte) {
	fmt.Println("Reading TSM Blocks...")
	// i := 1
	// for {
	// 	fmt.Println("***********************")
	// 	fmt.Println("Block", i)
	// 	fmt.Println("***********************")

	// 	// Get CRC!!
	// 	_, err = reader.Seek(int64(bOffset), io.SeekStart)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	var crcBuf = make([]byte, 4)
	// 	_, err = reader.Read(crcBuf)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	crcReader := bytes.NewReader(crcBuf)
	// 	var crc uint32
	// 	err = binary.Read(crcReader, binary.BigEndian, &crc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("crc ")
	// 	fmt.Println(crc)

	// 	var blockBuf = make([]byte, bSize-4)
	// 	_, err = reader.Read(blockBuf)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	crcVal := crc32.ChecksumIEEE(blockBuf)
	// 	fmt.Println("crcVal ")
	// 	fmt.Println(crcVal)

	// 	// currPos, err := reader.Seek(0, io.SeekCurrent)
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	// if currPos >= int64(footer.Pos) {
	// 	// 	break
	// 	// }
	// 	i = i + 1
	// }
}
