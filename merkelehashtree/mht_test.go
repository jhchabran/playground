package merkelehashtree

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	size := 1024 * 1024
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(i)
	}

	buf := bytes.NewReader(data)
	root, err := New(buf, 4096)
	if err != nil {
		t.Error(err)
	}

	res, _ := ReadAll(root)

	if len(data) != len(res) {
		t.Errorf("want len=%d got %d", len(data), len(res))
	}

	for i, c := range data {
		if c != res[i] {
			t.Errorf("want %v, got %v at index %d", c, res[i], i)
		}
	}
}
