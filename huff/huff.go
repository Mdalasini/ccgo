package main

import (
	"sort"
)

type HuffNode interface {
	isLeaf() bool
	getFreq() int
}

type HuffLeaf struct {
	char rune
	freq int
}

func (l HuffLeaf) isLeaf() bool {
	return true
}

func (l HuffLeaf) getFreq() int {
	return l.freq
}

type HuffBranch struct {
	freq  int
	left  HuffNode
	right HuffNode
}

func (n HuffBranch) isLeaf() bool {
	return false
}

func (n HuffBranch) getFreq() int {
	return n.freq
}

func sortNodesByFreq(nodes []HuffNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].getFreq() < nodes[j].getFreq()
	})
}

func insertSorted(nodes []HuffNode, node HuffNode) []HuffNode {
	idx := sort.Search(len(nodes), func(i int) bool {
		return nodes[i].getFreq() >= node.getFreq()
	})

	nodes = append(nodes, node)
	copy(nodes[idx+1:], nodes[idx:])
	nodes[idx] = node

	return nodes
}

func buildTree(nodes []HuffNode) HuffNode {
	sortNodesByFreq(nodes)
	for len(nodes) > 1 {
		left := nodes[0]
		right := nodes[1]
		merged := mergeNodes(left, right) // left should have lower freq
		nodes = nodes[2:]
		nodes = insertSorted(nodes, merged)
	}
	return nodes[0]
}

func mergeNodes(left, right HuffNode) HuffBranch {
	return HuffBranch{
		freq:  left.getFreq() + right.getFreq(),
		left:  left,
		right: right,
	}
}

func nodesFromFreq(frequencies map[rune]int) []HuffNode {
	nodes := make([]HuffNode, 0, len(frequencies))
	for char, freq := range frequencies {
		nodes = append(nodes, HuffLeaf{char: char, freq: freq})
	}
	return nodes
}
