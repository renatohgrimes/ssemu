package resources

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rasky/go-lzo"
)

type x4File struct {
	name     string
	contents []byte
}

func extractX4Files(hdPath string, resourcesDirPath string, cypherKey []byte) ([]x4File, error) {
	var err error
	hd, err := os.ReadFile(hdPath)
	if err != nil {
		return nil, err
	}
	hdreader := bytes.NewReader(hd)
	version, count := readHeader(hdreader)
	if version != 1 {
		return nil, errors.New("resource map is invalid")
	}
	var files []x4File
	var i int32
	for i = 0; i < count; i++ {
		resource, crc, size := readResourceHeader(hdreader, cypherKey)
		if strings.HasSuffix(resource, ".x4") {
			path := path.Join(resourcesDirPath, strconv.FormatUint(uint64(crc), 16))
			contents, err := readResourceContents(path, size, cypherKey)
			if err != nil {
				return nil, err
			}
			files = append(files, x4File{name: resource, contents: contents})
		}
	}
	return files, nil
}

func readHeader(reader *bytes.Reader) (version int32, count int32) {
	binary.Read(reader, binary.LittleEndian, &version)
	binary.Read(reader, binary.LittleEndian, &count)
	return
}

func readResourceHeader(reader *bytes.Reader, cypherKey []byte) (resourceName string, crc int64, size int32) {
	var headerSize int32
	binary.Read(reader, binary.LittleEndian, &headerSize)
	header := make([]byte, headerSize)
	binary.Read(reader, binary.LittleEndian, &header)
	decrypt(header, cypherKey, true)
	swapBytes(header)
	header = decompress(header, 272)
	decrypt(header, cypherKey, true)
	swapBytes(header)
	headerReader := bytes.NewReader(header)
	nameBuffer := make([]byte, 256)
	binary.Read(headerReader, binary.LittleEndian, &nameBuffer)
	binary.Read(headerReader, binary.LittleEndian, &crc)
	binary.Read(headerReader, binary.LittleEndian, &size)
	resourceName = string(trim(nameBuffer))[4:]
	return
}

func readResourceContents(path string, size int32, cypherKey []byte) ([]byte, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	swapBytes(contents)
	contents = decompress(contents, int(size))
	decrypt(contents, cypherKey, true)
	decrypt(contents, cypherKey, false)
	contents = decompress(contents, len(contents)*30)
	return contents, nil
}

func decrypt(buffer []byte, key []byte, capped bool) {
	size := len(buffer)
	if size > 256 && capped {
		size = 256
	}
	limit := 32
	if !capped {
		limit = 40
	}
	for i := 0; i < size; i++ {
		buffer[i] = ((buffer[i] >> 1) & 0x7F) | ((buffer[i] & 1) << 0x07)
		buffer[i] ^= key[(i % limit)]
	}
}

func swapBytes(buffer []byte) {
	length := len(buffer)
	capping := length
	if capping > 128 {
		capping = 128
	}
	for i := 0; i < (capping / 2); i++ {
		pos := length - 1 - i
		tmp := buffer[pos]
		buffer[pos] = buffer[i]
		buffer[i] = tmp
	}
}

func decompress(buffer []byte, realSize int) []byte {
	b, _ := lzo.Decompress1X(bytes.NewReader(buffer), len(buffer), realSize)
	return b
}

func trim(buffer []byte) []byte {
	length := len(buffer)
	var pos int
	for pos = 0; buffer[pos] != 0x00 && pos < length; pos++ {
	}
	return buffer[0:pos]
}
