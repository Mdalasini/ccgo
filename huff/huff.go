package main

import (
	"container/heap"
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

// BuildCodeTable walks the Huffman tree and returns a map from rune to its
// variable-length prefix code (as a string of '0'/'1').
func BuildCodeTable(root *HuffNode) map[rune]string {
	codes := make(map[rune]string)
	if root == nil {
		return codes
	}
	if root.isLeaf() {
		codes[root.Char] = "0"
		return codes
	}
	var walk func(n *HuffNode, path string)
	walk = func(n *HuffNode, path string) {
		if n.isLeaf() {
			codes[n.Char] = path
			return
		}
		if n.Left != nil {
			walk(n.Left, path+"0")
		}
		if n.Right != nil {
			walk(n.Right, path+"1")
		}
	}
	walk(root, "")
	return codes
}
