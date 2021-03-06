package vpk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

const (
	signature = 0x55aa1234

	hSizeV1 = 0x0d
	hSizeV2 = 0x1b

	dirName = "_dir.vpk"
)

var (
	order       = binary.LittleEndian
	ErrBroken   = errors.New("This vpk file(s) is ErrBroke.")
	ErrChecksum = errors.New("VPK checksum error.")
)

type vpkHeader struct {
	Version      uint32
	TreeLength   uint32
	FooterLength uint32
}

type FileHeader struct {
	CRC   uint32
	entry *fileEntry
}

type fileEntry struct {
	PreloadBytes uint16
	ArchiveIndex uint16
	EntryOffset  uint32
	EntryLength  uint32
}

type File struct {
	FileHeader
	Name   string
	parent *ReadCloser
}

type Reader struct {
	File []*File
}

type ReadCloser struct {
	Reader
	head *vpkHeader
	fd   *os.File
	fs   map[int]*os.File
}

func (r *ReadCloser) Close() {
	r.fd.Close()
	for _, f := range r.fs {
		f.Close()
	}
}

//OpenReader is open vpk file function
//if u want extract...
//	Solo file: ("path/to/file.vpk")
//	Sequence files: ("path/to/file_dir.vpk")
//	Directory: ("path/to/")
func OpenReader(path string) (*ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	r := new(ReadCloser)
	if err := r.init(f, fi.Size()); err != nil {
		f.Close()
		return nil, err
	}
	r.fd = f

	if strings.HasSuffix(f.Name(), dirName) {
		r.fs = make(map[int]*os.File)
		for i := 1; i < 1000; i++ {
			spath := fmt.Sprintf("%s%03d.vpk", path[:len(path)-7], i)
			sf, err := os.Open(spath)
			if err != nil {
				break
			}
			r.fs[i] = sf
		}
	}

	return r, nil
}

func (f *File) Open() (io.ReadCloser, error) {
	rc := new(fileReader)
	p := f.parent
	if f.entry.ArchiveIndex == 0x7fff {
		rc.f = p.fd
		rc.offset = int64(f.entry.EntryOffset + p.head.TreeLength)
	} else {
		rc.f = p.fs[int(f.entry.ArchiveIndex)]
		if rc.f != nil {
			return nil, fmt.Errorf("Not found archive file.")
		}
		rc.offset = int64(f.entry.EntryOffset)
	}
	switch f.parent.head.Version {
	case 1:
		rc.offset += hSizeV1
	case 2:
		rc.offset += hSizeV2
	}
	rc.offset++
	rc.length = int64(f.entry.EntryLength)
	rc.hash = crc32.NewIEEE()
	rc.crc32 = f.CRC
	return rc, nil
}

func (z *ReadCloser) init(r io.Reader, size int64) error {
	head, err := readHeader(r, size)
	if err != nil {
		return err
	}
	z.head = head
	files, err := readDirectory(r, size)
	if err != nil {
		return err
	}
	for _, file := range files {
		file.parent = z
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
	header := new(FileHeader)
	binary.Read(r, order, &header.CRC)
	entry := new(fileEntry)
	binary.Read(r, order, entry)
	header.entry = entry
	var terminator uint16
	binary.Read(r, order, &terminator)
	if terminator != 0xffff {
		return nil, ErrBroken
	}
	return header, nil
}

type fileReader struct {
	io.ReadCloser
	hash   hash.Hash32
	f      *os.File
	offset int64
	length int64
	closed bool
	crc32  uint32
}

func (r *fileReader) Read(b []byte) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("Closed.")
	}
	if r.length == 0 {
		return 0, io.EOF
	}
	_, err := r.f.Seek(r.offset, 0)
	if err != nil {
		return 0, err
	}
	size := len(b)
	l := int(r.length)
	if size > l {
		size = l
	}
	s, err := r.f.Read(b[:size])
	r.hash.Write(b[:size])
	r.offset += int64(s)
	r.length -= int64(s)
	if err == io.EOF {
		if r.hash.Sum32() != r.crc32 {
			err = ErrChecksum
		}
	}
	return s, err
}

func (r *fileReader) Close() error {
	r.closed = true
	return nil
}

func (r *Reader) FindFile(name string) *File {
	for _, file := range r.File {
		if name == file.Name {
			return file
		}
	}
	return nil
}
