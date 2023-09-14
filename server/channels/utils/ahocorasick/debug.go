package ahocorasick

import "fmt"

func (n *Node) DEBUG__CountNodes() int {
	count := 1
	for _, node := range n.Edges {
		if node == n {
			continue
		}
		count += node.DEBUG__CountNodes()
	}

	return count
}

func (n *Node) DEBUG__PrintStrings(prefix string) {
	if n.Leaf.HasValues() {
		fmt.Println(prefix)
	}
	for k, node := range n.Edges {
		if node == n {
			continue
		}
		node.DEBUG__PrintStrings(prefix + k)
	}
}
