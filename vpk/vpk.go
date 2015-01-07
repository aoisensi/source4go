package vpk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	signature = 0x55aa1234
)

var (
	order  = binary.LittleEndian
	broken = errors.New("This vpk file(s) is broken.")
)

type vpkHeader struct {
	Version      uint32
	TreeLength   uint32
	FooterLength uint32
}

type FileHeader struct {
	CRC          uint32
	PreloadBytes uint16
	ArchiveIndex uint16
	EntryOffset  uint32
	EntryLength  uint32
}

type File struct {
	FileHeader
	Name string
}

type Reader struct {
	File []*File
}

type ReadCloser struct {
	f *os.File
	Reader
}

func (r *ReadCloser) Close() {
	r.f.Close()
}

//OpenReader is open vpk file function
//if u want extract...
//	Solo file: ("path/to/file.vpk")
//	Sequence files: ("path/to/file_dir.vpk")
//	Sequence files: ("path/to/file_dir.vpk", "path/to/file_001.vpk")
//	Directory: ("path/to/")
func OpenReader(path string, more ...string) (*ReadCloser, error) {
	r := new(ReadCloser)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if err := r.init(f, fi.Size()); err != nil {
		f.Close()
		return nil, err
	}
	r.f = f
	return r, nil
}

func (z *Reader) init(r io.Reader, size int64) error {
	_, err := readHeader(r, size)
	if err != nil {
		return err
	}
	files, err := readDirectory(r, size)
	if err != nil {
		return err
	}
	z.File = files
	return nil
}

func readHeader(r io.Reader, size int64) (*vpkHeader, error) {
	var sign uint32
	binary.Read(r, order, &sign)
	if sign != signature {
		return nil, fmt.Errorf("The file is not vpk file.")
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

func readString(r io.Reader) string {
	buf := make([]byte, 0, 0xff)
	var f byte
	for {
		binary.Read(r, order, &f)
		if f == 0 {
			break
		}
		buf = append(buf, f)
	}
	return string(buf)
}

func readDirectory(r io.Reader, size int64) ([]*File, error) {
	files := []*File{}
	for {
		ext := readString(r)
		if ext == "" {
			break
		}
		for {
			path := readString(r)
			if path == "" {
				break
			}
			for {
				name := readString(r)
				if name == "" {
					break
				}
				entry, err := readFileInfo(r, size)
				if err != nil {
					return nil, err
				}
				file := new(File)
				file.FileHeader = *entry
				file.Name = path + "/" + name + "." + ext
				files = append(files, file)
			}
		}
	}
	return files, nil
}

func readFileInfo(r io.Reader, size int64) (*FileHeader, error) {
	entry := new(FileHeader)
	binary.Read(r, order, entry)
	var terminator uint16
	binary.Read(r, order, &terminator)
	if terminator != 0xffff {
		return nil, broken
	}
	return entry, nil
}
