// Package infinity
//
// http://oeis.org/A004718
package infinity

import "github.com/bh90210/mlsic/graph"

func Series(n, offset int) []int {
	var s []int
	n = n * 2
	for i := 0; i <= (n + 1); i++ {
		switch {
		case i == 0:
			s = append(s, offset)

		case i%2 == 1:
			node := i - 1
			node = node / 2
			s = append(s, s[node]+1)

		case i%2 == 0:
			node := i / 2
			s = append(s, s[node]*-1)

		}
	}
	return s
}

func BinaryTree(l int) *graph.Node {
	// We start constructing the tree from level 1 upwards.
	n1 := &graph.Node{Key: 1}
	// Purpose of slice s(eries) is to hold the previous level's graph.Nodes.
	s := []*graph.Node{n1}
	// Init rule a(0) = 0 .
	n := &graph.Node{Key: 0, Right: n1}

	// Start looping from level 1 to level l.
	// Since init rule covers zero level we start constructing from level 1 upwards.
	for i := 1; i < l; i++ {
		// Helper slices temporarily holding newly
		// created nodes until level cycle completes.
		var h []*graph.Node
		// Range through previous level's nodes.
		for _, v := range s {
			// a(2n) = -a(n)
			v.Left = &graph.Node{Key: v.Key * -1}
			// a(2n+1) = a(n) + 1
			v.Right = &graph.Node{Key: v.Key + 1}
			h = append(h, v.Left, v.Right)
			// If node's key equals i means level reached the end
			// and we need to break and move to the next.
			if v.Key == i {
				// Assign this level's newly created slice of nodes h to s.
				s = h
				break
			}
		}
	}

	return n
}
