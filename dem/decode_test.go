package dem

import (
	"os"
	"testing"

	"github.com/k0kubun/pp"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("test.dem")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	demo, err := NewDemo(f)
	if err != nil {
		t.Fatal(err)
	}
	pp.Println(demo.Header())
}
