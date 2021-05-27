package mht

import (
	"bytes"
	"crypto/sha256"
)

type Node struct {
	left  *Node
	right *Node
	block *block
	Hash  []byte
}

type block struct {
	index int64
	data  []byte
}

func New(r *bytes.Reader, chunkSize int64) (*Node, error) {
	chunkCount := r.Size() / chunkSize
	rem := r.Size() % chunkSize
	if rem > 0 {
		chunkCount += 1
	}
	leaves := make([]*Node, 0, chunkCount)
	b := make([]byte, chunkSize)

	for i := int64(0); i < chunkCount; i++ {
		c, err := r.Read(b)
		if err != nil {
			return nil, err
		}
		n := &Node{
			block: &block{
				index: i,
				data:  make([]byte, chunkSize),
			},
		}
		copy(n.block.data, b[0:c])
		n.block.data = n.block.data[0:c]
		h := sha256.New()
		_, err = h.Write(n.block.data)
		if err != nil {
			return nil, err
		}
		n.Hash = h.Sum(nil)
		leaves = append(leaves, n)
	}

	nodes := leaves
	for {
		if len(nodes) == 1 {
			return nodes[0], nil
		}
		above := []*Node{}
		for i := 0; i < len(nodes); {
			n := &Node{}
			h := sha256.New()
			if i+1 < len(nodes) {
				n.left = nodes[i]
				n.right = nodes[i+1]
				_, err := h.Write(n.left.Hash)
				if err != nil {
					return nil, err
				}
				n.Hash = h.Sum(n.right.Hash)
				i += 2
			} else {
				n.left = nodes[i]
				_, err := h.Write(n.left.Hash)
				if err != nil {
					return nil, err
				}
				n.Hash = h.Sum(nil)
				i++
			}
			above = append(above, n)
		}
		nodes = above
	}
}

func ReadAll(root *Node) ([]byte, error) {
	nodes := []*Node{root}
	data := []byte{}
	for {
		if len(nodes) == 0 {
			return data, nil
		}
		below := []*Node{}
		for _, n := range nodes {
			if n.block != nil {
				data = append(data, n.block.data...)
				continue
			}
			if n.left != nil {
				below = append(below, n.left)
			}
			if n.right != nil {
				below = append(below, n.right)
			}
		}
		nodes = below
	}
}
