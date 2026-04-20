package main

import "testing"

func newHuffLeaf(t *testing.T, char rune, freq int) HuffLeaf {
	return HuffLeaf{
		char: char,
		freq: freq,
	}
}

// Tree traversal and property helpers

func countLeaves(node HuffNode) int {
	if node.isLeaf() {
		return 1
	}
	branch := node.(HuffBranch)
	return countLeaves(branch.left) + countLeaves(branch.right)
}

func sumLeafFreqs(node HuffNode) int {
	if node.isLeaf() {
		return node.getFreq()
	}
	branch := node.(HuffBranch)
	return sumLeafFreqs(branch.left) + sumLeafFreqs(branch.right)
}

func collectLeaves(node HuffNode) []rune {
	if node.isLeaf() {
		return []rune{node.(HuffLeaf).char}
	}
	branch := node.(HuffBranch)
	return append(collectLeaves(branch.left), collectLeaves(branch.right)...)
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

// verifyExactStructure checks structure for simple trees with two leaves.
func verifyExactStructure(t *testing.T, tree HuffNode, expectedRoot int, leftChar, rightChar rune) {
	t.Helper()
	if tree.getFreq() != expectedRoot {
		t.Errorf("root freq: got %d, want %d", tree.getFreq(), expectedRoot)
	}
	branch := tree.(HuffBranch)
	if leftChar != 0 {
		if branch.left.(HuffLeaf).char != leftChar {
			t.Errorf("left char: got %q, want %q", branch.left.(HuffLeaf).char, leftChar)
		}
	}
	if rightChar != 0 {
		if branch.right.(HuffLeaf).char != rightChar {
			t.Errorf("right char: got %q, want %q", branch.right.(HuffLeaf).char, rightChar)
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
