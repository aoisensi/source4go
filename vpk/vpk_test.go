package vpk

import (
	"bytes"
	"io"
	"testing"
)

func TestFileCheck(t *testing.T) {
	r, err := OpenReader("test.vpk")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	if len(r.File) != 1 {
		panic(err)
	}
	file := r.File[0]
	name := file.Name
	if name != "the/cake/is/a/lie/glados.txt" {
		t.Error(name)
	}
	rc, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer rc.Close()
	buf := make([]byte, 0x20)
	_, err = rc.Read(buf)
	if err != io.EOF && err != nil {
		t.Error(err)
	}
	text := []byte("the cake is a lie...")

	if !bytes.Equal(buf[:len(text)], text) {
		t.Error("Failed to read file")
	}
}
