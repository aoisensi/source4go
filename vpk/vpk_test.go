package vpk

import "testing"

func TestFileCheck(t *testing.T) {
	r, err := OpenReader("test.vpk")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	if len(r.File) != 1 {
		panic(err)
	}
	name := r.File[0].Name
	if name != "the/cake/is/a/lie/glados.txt" {
		t.Error(name)
	}
}
