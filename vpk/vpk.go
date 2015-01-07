package vpk

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

const (
	signature = 0x55aa1234
)

var (
	order = binary.BigEndian
)

type vpkHeader struct {
	Version      uint32
	TreeLength   uint32
	FooterLength uint32
}

type FileHeader struct {
	Name string

	CreatorVersion     uint16
	ReaderVersion      uint16
	Flags              uint16
	Method             uint16
	ModifiedTime       uint16 // MS-DOS time
	ModifiedDate       uint16 // MS-DOS date
	CRC32              uint32
	CompressedSize     uint32 // deprecated; use CompressedSize64
	UncompressedSize   uint32 // deprecated; use UncompressedSize64
	CompressedSize64   uint64
	UncompressedSize64 uint64
	Extra              []byte
	ExternalAttrs      uint32 // Meaning depends on CreatorVersion
	Comment            string
}

type File struct {
	FileHeader
}

type Reader struct {
	File []*File
}

type ReadCloser struct {
	f *os.File
	Reader
}

//OpenReader is open vpk file function
//if u want extract...
//	Solo file: ("path/to/file.vpk")
//	Sequence files: ("path/to/file_dir.vpk")
//	Sequence files: ("path/to/file_dir.vpk", "path/to/file_001.vpk")
//	Directory: ("path/to/")
func OpenReader(path string, more ...string) (*ReadCloser, error) {
	if len(more) == 0 {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {

		}
		r := new(ReadCloser)
		if err := r.init(f, fi.Size()); err != nil {
			f.Close()
			return nil, err
		}
		r.f = f
	}
	return nil, nil
}

func (z *Reader) init(r io.Reader, size int64) error {
	_, err := readHeader(r, size)
	if err != nil {
		return err
	}
	return nil
}

func readHeader(r io.Reader, size int64) (*vpkHeader, error) {
	var sign uint32
	binary.Read(r, order, &sign)
	if sign != signature {
		return nil, errors.New("The file is not vpk file.")
	}

	head := new(vpkHeader)
	binary.Read(r, order, &head.Version)
	if head.Version != 1 && head.Version != 2 {
		return nil, errors.New("The vpk file version is not support")
	}

	binary.Read(r, order, &head.TreeLength)
	if head.Version == 2 {
		var dummy uint32
		binary.Read(r, order, &dummy)
		binary.Read(r, order, &head.FooterLength)
		binary.Read(r, order, &dummy)
		binary.Read(r, order, &dummy)
	}

	return head, nil
}
