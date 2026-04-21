package main

import (
	"fmt"
	"log"
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
	freq      int
	leftNode  HuffNode
	rightNode HuffNode
}

func (n HuffBranch) isLeaf() bool {
	return false
}

func (n HuffBranch) getFreq() int {
	return n.freq
}

func (n HuffBranch) genCodes(path string, pathMap *map[rune]string) {
	visitNode := func(node HuffNode, code int) {
		newPath := fmt.Sprintf("%s%d", path, code)
		switch node := node.(type) {
		case HuffBranch:
			node.genCodes(newPath, pathMap)
		case HuffLeaf:
			if _, ok := (*pathMap)[node.char]; ok {
				log.Fatalf("%c already exists in map", node.char)
			}
			(*pathMap)[node.char] = newPath
		default:
			log.Fatalf("unknown type for node: %T\n", node)
		}
	}

	visitNode(n.leftNode, 0)
	visitNode(n.rightNode, 1)
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

func buildTree(nodes []HuffNode) (HuffBranch, error) {
	if len(nodes) == 0 {
		return HuffBranch{}, fmt.Errorf("cannot build tree from empty node list")
	}

	if len(nodes) == 1 {
		if leaf, ok := nodes[0].(HuffLeaf); ok {
			return HuffBranch{leftNode: leaf, freq: leaf.freq}, nil
		}
		return nodes[0].(HuffBranch), nil
	}

	sortNodesByFreq(nodes)
	for len(nodes) > 1 {
		left := nodes[0]
		right := nodes[1]
		merged := mergeNodes(left, right) // left should have lower freq
		nodes = nodes[2:]
		nodes = insertSorted(nodes, merged) // maintains sort
	}
	branch, ok := nodes[0].(HuffBranch)
	if !ok {
		return HuffBranch{}, fmt.Errorf("unexpected node type after merge: %T", nodes[0])
	}

	return branch, nil
}

func mergeNodes(left, right HuffNode) HuffBranch {
	return HuffBranch{
		freq:      left.getFreq() + right.getFreq(),
		leftNode:  left,
		rightNode: right,
	}
}

func nodesFromFreq(frequencies map[rune]int) []HuffNode {
	nodes := make([]HuffNode, 0, len(frequencies))
	for char, freq := range frequencies {
		nodes = append(nodes, HuffLeaf{char: char, freq: freq})
	}
	return nodes
}
