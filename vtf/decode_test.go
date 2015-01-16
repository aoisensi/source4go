package vtf

import (
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("test.vtf")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	_, err = Decode(f)
	if err != nil {
		t.Error(err)
	}
}
