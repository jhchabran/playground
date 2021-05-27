package mht_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/jhchabran/playground/mht"
)

func newBlob(size int64) []byte {
	data := make([]byte, size)
	for i := int64(0); i < size; i++ {
		data[i] = byte(i % 255)
	}
	return data
}

func TestMerkeleHashTree(t *testing.T) {
	blob := newBlob(1024 * 1024)
	r := bytes.NewReader(blob)

	var root *mht.Node
	var err error
	t.Run("New", func(t *testing.T) {
		root, err = mht.New(r, 4096)
		if err != nil {
			t.Error(err)
		}

		if root.Hash == nil {
			t.Errorf("want not nil root Hash, got nil")
		}
	})

	t.Run("ReadAll", func(t *testing.T) {
		res, _ := mht.ReadAll(root)

		if len(blob) != len(res) {
			t.Errorf("want len=%d got %d", len(blob), len(res))
		}

		for i, c := range blob {
			if c != res[i] {
				t.Errorf("want %v, got %v at index %d", c, res[i], i)
			}
		}
	})

	t.Run("Hash differs on diff data", func(t *testing.T) {
		otherBlob := newBlob(1024 * 1024)
		otherBlob[128] = 0x00

		buf := bytes.NewReader(otherBlob)
		otherRoot, err := mht.New(buf, 4096)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(root.Hash, otherRoot.Hash) {
			t.Errorf("want Hashes to differ, but got the same value")
		}
	})
}

func BenchmarkNew(b *testing.B) {
	for size := int64(256); float64(size) <= math.Pow(256, 3); size *= 256 {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			blob := newBlob(size)
			b.ResetTimer()
			b.StopTimer()

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				_, _ = mht.New(bytes.NewReader(blob), 256)
			}
		})
	}
}

func BenchmarkReadAll(b *testing.B) {
	for size := int64(256); float64(size) <= math.Pow(256, 3); size *= 256 {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			blob := newBlob(size)
			root, _ := mht.New(bytes.NewReader(blob), 256)
			b.ResetTimer()
			b.StopTimer()

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				_, _ = mht.ReadAll(root)
			}
		})
	}
}
