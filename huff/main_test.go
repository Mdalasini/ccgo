package main

import "testing"

// Tree traversal and property helpers

func countLeaves(node HuffNode) int {
	if node.isLeaf() {
		return 1
	}
	branch := node.(HuffBranch)
	return countLeaves(branch.leftNode) + countLeaves(branch.rightNode)
}

func sumLeafFreqs(node HuffNode) int {
	if node.isLeaf() {
		return node.getFreq()
	}
	branch := node.(HuffBranch)
	return sumLeafFreqs(branch.leftNode) + sumLeafFreqs(branch.rightNode)
}

func collectLeaves(node HuffNode) []rune {
	if node.isLeaf() {
		return []rune{node.(HuffLeaf).char}
	}
	branch := node.(HuffBranch)
	return append(collectLeaves(branch.leftNode), collectLeaves(branch.rightNode)...)
}

func verifyTreeInvariants(t *testing.T, tree HuffNode, expectedLeaves []HuffLeaf) {
	t.Helper()

	// Property 1: leaf count matches input
	leafCount := countLeaves(tree)
	if leafCount != len(expectedLeaves) {
		t.Errorf("leaf count mismatch: got %d, want %d", leafCount, len(expectedLeaves))
	}

	// Property 2: sum of leaf frequencies equals root frequency
	sum := sumLeafFreqs(tree)
	if sum != tree.getFreq() {
		t.Errorf("frequency sum mismatch: sum=%d, root=%d", sum, tree.getFreq())
	}

	// Property 3: all expected chars are present
	chars := collectLeaves(tree)
	charSet := make(map[rune]bool)
	for _, c := range chars {
		charSet[c] = true
	}
	for _, leaf := range expectedLeaves {
		if !charSet[leaf.char] {
			t.Errorf("expected char %q not found in tree", leaf.char)
		}
	}
}

func TestBuildTree(t *testing.T) {
	tests := []struct {
		name     string
		leaves   []HuffLeaf
		rootFreq int
	}{
		{
			name:     "two nodes",
			leaves:   []HuffLeaf{{char: 'A', freq: 10}, {char: 'B', freq: 20}},
			rootFreq: 30,
		},
		{
			name:     "three nodes",
			leaves:   []HuffLeaf{{char: 'A', freq: 5}, {char: 'B', freq: 10}, {char: 'C', freq: 15}},
			rootFreq: 30,
		},
		{
			name:     "single node",
			leaves:   []HuffLeaf{{char: 'X', freq: 42}},
			rootFreq: 42,
		},
		{
			name:     "complex tree",
			leaves:   []HuffLeaf{{char: 'C', freq: 32}, {char: 'D', freq: 42}, {char: 'E', freq: 120}, {char: 'K', freq: 7}, {char: 'L', freq: 42}, {char: 'M', freq: 24}, {char: 'U', freq: 37}, {char: 'Z', freq: 2}},
			rootFreq: 306,
		},
		{
			name:     "all same frequency",
			leaves:   []HuffLeaf{{char: 'A', freq: 10}, {char: 'B', freq: 10}, {char: 'C', freq: 10}},
			rootFreq: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := make([]HuffNode, len(tt.leaves))
			for i, leaf := range tt.leaves {
				nodes[i] = leaf
			}
			tree := buildTree(nodes)

			// Always verify invariants
			verifyTreeInvariants(t, tree, tt.leaves)

			// Verify root frequency
			if tree.getFreq() != tt.rootFreq {
				t.Errorf("root freq: got %d, want %d", tree.getFreq(), tt.rootFreq)
			}
		})
	}
}

func TestGenCodesComplexTree(t *testing.T) {
	leaves := []HuffLeaf{
		{char: 'C', freq: 32},
		{char: 'D', freq: 42},
		{char: 'E', freq: 120},
		{char: 'K', freq: 7},
		{char: 'L', freq: 42},
		{char: 'M', freq: 24},
		{char: 'U', freq: 37},
		{char: 'Z', freq: 2},
	}
	nodes := make([]HuffNode, len(leaves))
	for i, leaf := range leaves {
		nodes[i] = leaf
	}
	tree := buildTree(nodes)

	initialPath := ""
	pathMap := make(map[rune]string)
	tree.(HuffBranch).genCodes(initialPath, &pathMap)

	expected := map[rune]string{
		'C': "1110",
		'D': "101",
		'E': "0",
		'K': "111101",
		'L': "110",
		'M': "11111",
		'U': "100",
		'Z': "111100",
	}

	for char, want := range expected {
		got, ok := pathMap[char]
		if !ok {
			t.Errorf("char %q not found in pathMap", char)
			continue
		}
		if got != want {
			t.Errorf("code for %q: got %s, want %s", char, got, want)
		}
	}

	if len(pathMap) != len(expected) {
		t.Errorf("pathMap size: got %d, want %d", len(pathMap), len(expected))
	}
}
