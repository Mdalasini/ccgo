package main

import (
	"container/heap"
	"fmt"
)

// HuffNode represents a node in a Huffman tree.
// Leaf nodes store a rune; internal nodes store their left and right children.
type HuffNode struct {
	Char  rune
	Freq  int
	Left  *HuffNode
	Right *HuffNode
}

// isLeaf returns true if this node has no children (i.e. it represents a single rune).
func (n *HuffNode) isLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// HuffHeap implements a min-heap of *HuffNode ordered by frequency.
type HuffHeap []*HuffNode

func (h HuffHeap) Len() int { return len(h) }

func (h HuffHeap) Less(i, j int) bool {
	if h[i].Freq == h[j].Freq {
		return h[i].Char < h[j].Char
	}
	return h[i].Freq < h[j].Freq
}

func (h HuffHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *HuffHeap) Push(x any) {
	*h = append(*h, x.(*HuffNode))
}

func (h *HuffHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[:n-1]
	return x
}

// BuildHuffTree builds a Huffman tree from a frequency map.
// It returns the root of the tree, or nil if the frequency map is empty.
func BuildHuffTree(freq map[rune]int) *HuffNode {
	if len(freq) == 0 {
		return nil
	}

	h := &HuffHeap{}
	heap.Init(h)
	for r, f := range freq {
		heap.Push(h, &HuffNode{Char: r, Freq: f})
	}

	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffNode)
		right := heap.Pop(h).(*HuffNode)

		parent := &HuffNode{
			Char:  0, // internal nodes have no meaningful rune
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}

		heap.Push(h, parent)
	}

	return heap.Pop(h).(*HuffNode)
}

func (n *HuffNode) Traverse(path string) (rune, error) {
	cur := n
	for i, bit := range path {
		if cur.isLeaf() {
			return 0, fmt.Errorf("path too long: landed on leaf at offset %d", i)
		}
		switch bit {
		case '0':
			cur = cur.Left
		case '1':
			cur = cur.Right
		default:
			return 0, fmt.Errorf("invalid path character %q at offset %d", bit, i)
		}
		if cur == nil {
			return 0, fmt.Errorf("path leads to nil node at offset %d", i)
		}
	}
	if !cur.isLeaf() {
		return 0, fmt.Errorf("path ended at internal node")
	}
	return cur.Char, nil
}
