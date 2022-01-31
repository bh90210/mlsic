package graph

type Node struct {
	Key   int
	Left  *Node
	Right *Node
}

func Breadth(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		result = append(result, currentNode.Key)
		visit = visit[1:]
		if currentNode.Left != nil {
			visit = append(visit, currentNode.Left)
		}
		if currentNode.Right != nil {
			visit = append(visit, currentNode.Right)
		}
	}
	return result
}

func DepthPreOrder(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		// Append node's value to result slice.
		result = append(result, currentNode.Key)
		if currentNode.Right != nil {
			// Makes slices one position bigger.
			visit = append(visit, &Node{})
			// Copy everything starting from position 1
			copy(visit[2:], visit[1:])
			// Assign to visit[1] the next right.
			visit[1] = currentNode.Right
		}
		// If node has left child append it to visit slice.
		if currentNode.Left != nil {
			visit[0] = currentNode.Left
			// Pop node from visit slice.
		} else if currentNode.Left == nil {
			visit = visit[1:]
		}
	}

	return result
}

func ReverseDepthPreOrder(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		// Append node's value to result slice.
		result = append(result, currentNode.Key)
		if currentNode.Left != nil {
			// Makes slices one position bigger.
			visit = append(visit, &Node{})
			// Copy everything starting from position 1
			copy(visit[2:], visit[1:])
			// Assign to visit[1] the next right.
			visit[1] = currentNode.Left
		}
		// If node has left child append it to visit slice.
		if currentNode.Right != nil {
			visit[0] = currentNode.Right
			// Pop node from visit slice.
		} else if currentNode.Right == nil {
			visit = visit[1:]
		}
	}

	return result
}

func DepthInOrder(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		if currentNode.Left != nil {
			visit = append(visit, &Node{})
			copy(visit[1:], visit[0:])
			visit[0] = visit[1].Left
			visit[1].Left = nil
		} else if currentNode.Left == nil {
			result = append(result, currentNode.Key)
			if currentNode.Right != nil {
				visit[0] = currentNode.Right
			} else {
				visit = visit[1:]
			}
		}
	}

	return result
}

func ReverseDepthInOrder(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		if currentNode.Right != nil {
			visit = append(visit, &Node{})
			copy(visit[1:], visit[0:])
			visit[0] = visit[1].Right
			visit[1].Right = nil
		} else if currentNode.Right == nil {
			result = append(result, currentNode.Key)
			if currentNode.Left != nil {
				visit[0] = currentNode.Left
			} else {
				visit = visit[1:]
			}
		}
	}

	return result
}

func DepthPostOrder(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		if currentNode.Left != nil {
			visit = append(visit, &Node{})
			copy(visit[1:], visit[0:])
			visit[0] = visit[1].Left
			visit[1].Left = nil
		} else if currentNode.Left == nil {
			if currentNode.Right != nil {
				// visit[0] = currentNode.right
				visit = append(visit, &Node{})
				copy(visit[1:], visit[0:])
				visit[0] = visit[1].Right
				visit[1].Right = nil
			} else if currentNode.Right == nil {
				result = append(result, currentNode.Key)
				visit = visit[1:]
			}
		}
	}

	return result
}

func ReverseDepthPostOrder(n *Node) []int {
	var result []int
	visit := []*Node{n}
	for len(visit) > 0 {
		currentNode := visit[0]
		if currentNode.Right != nil {
			visit = append(visit, &Node{})
			copy(visit[1:], visit[0:])
			visit[0] = visit[1].Right
			visit[1].Right = nil
		} else if currentNode.Right == nil {
			if currentNode.Left != nil {
				// visit[0] = currentNode.left
				visit = append(visit, &Node{})
				copy(visit[1:], visit[0:])
				visit[0] = visit[1].Left
				visit[1].Left = nil
			} else if currentNode.Left == nil {
				result = append(result, currentNode.Key)
				visit = visit[1:]
			}
		}
	}

	return result
}
