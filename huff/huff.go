package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"unicode/utf8"
)

type HuffNode interface {
	getFreq() int
	genCodes(path []byte, pathMap map[rune][]byte)
}

type HuffLeaf struct {
	char rune
	freq int
}

func (l HuffLeaf) getFreq() int {
	return l.freq
}

func (l HuffLeaf) genCodes(path []byte, pathMap map[rune][]byte) {
	pathCopy := make([]byte, len(path))
	copy(pathCopy, path) // avoid corruption by siblings
	pathMap[l.char] = pathCopy
}

type HuffBranch struct {
	freq      int
	leftNode  HuffNode
	rightNode HuffNode
}

func (b HuffBranch) getFreq() int {
	return b.freq
}

func (b HuffBranch) genCodes(path []byte, pathMap map[rune][]byte) {
	if b.leftNode != nil {
		b.leftNode.genCodes(append(path, '0'), pathMap)
	}
	if b.rightNode != nil {
		b.rightNode.genCodes(append(path, '1'), pathMap)
	}
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

func countCharFrequencies(r io.Reader) map[rune]int {
	frequencies := make(map[rune]int)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		r, _ := utf8.DecodeRuneInString(scanner.Text())
		frequencies[r]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	return frequencies
}

func nodesFromFreq(frequencies map[rune]int) []HuffNode {
	nodes := make([]HuffNode, 0, len(frequencies))
	for char, freq := range frequencies {
		nodes = append(nodes, HuffLeaf{char: char, freq: freq})
	}
	return nodes
}
