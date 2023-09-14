package ahocorasick

import (
	"strings"
)

type Node struct {
	Edges map[string]*Node
	Leaf  *Leaf
}

func NewNode() *Node {
	return &Node{
		Edges: map[string]*Node{},
		Leaf:  &Leaf{},
	}
}

type KeywordString []KeywordTerm

type KeywordTerm struct {
	Type            string
	Term            string
	CaseInsensitive bool
}

func (n *Node) AddString(keywordTerms KeywordString, value LeafObject) {
	n.addString(keywordTerms, value, "")
}

func handleCharacterAdd(n *Node, edgeCharacter string, prefix string, rest string, restKeywordTerms KeywordString, value LeafObject, caseInsensitive bool) {
	node, ok := n.Edges[edgeCharacter]
	if !ok {
		node = NewNode()
		node.Leaf.Term = prefix + edgeCharacter
		n.Edges[edgeCharacter] = node
	}

	node.addString(append(KeywordString{{Type: "string", Term: rest, CaseInsensitive: caseInsensitive}}, restKeywordTerms...), value, prefix+edgeCharacter)
}

func (n *Node) addString(keywordTerms KeywordString, value LeafObject, prefix string) {
	if len(keywordTerms) == 0 {
		n.Leaf.AddValue(value)
		return
	}

	currentKeywordTerm := keywordTerms[0]
	restKeywordTerms := keywordTerms[1:]

	switch currentKeywordTerm.Type {
	case "string":
		if len(currentKeywordTerm.Term) == 0 {
			n.addString(restKeywordTerms, value, prefix)
			return
		}

		edgeCharacter := currentKeywordTerm.Term[:1]
		rest := currentKeywordTerm.Term[1:]
		if currentKeywordTerm.CaseInsensitive {
			lowercaseTerm := strings.ToLower(edgeCharacter)
			uppercaseTerm := strings.ToUpper(edgeCharacter)
			if lowercaseTerm != uppercaseTerm {
				handleCharacterAdd(n, "x"+lowercaseTerm, prefix, rest, restKeywordTerms, value, true)
			} else {
				handleCharacterAdd(n, lowercaseTerm, prefix, rest, restKeywordTerms, value, true)
			}
		} else {
			handleCharacterAdd(n, edgeCharacter, prefix, rest, restKeywordTerms, value, false)
		}

	case "alphanumeric":
		fallthrough
	case "wordseparator":
		newPrefix := prefix + "{{" + currentKeywordTerm.Type + "}}"
		node, ok := n.Edges[currentKeywordTerm.Type]
		if !ok {
			node = NewNode()
			node.Leaf.Term = newPrefix
			n.Edges[currentKeywordTerm.Type] = node
		}

		node.addString(restKeywordTerms, value, newPrefix)
	case "alphanumeric*":
		fallthrough
	case "wordseparator*":
		newPrefix := prefix + "{{" + currentKeywordTerm.Type + "}}"
		node, ok := n.Edges[currentKeywordTerm.Type]
		if !ok {
			node = NewNode()
			node.Leaf.Term = newPrefix
			node.Edges[currentKeywordTerm.Type] = node
			n.Edges[currentKeywordTerm.Type] = node
		}

		node.addString(restKeywordTerms, value, newPrefix)
	}

}

func (n *Node) RemoveString(keywordTerms KeywordString, value LeafObject) {
	if len(keywordTerms) == 0 {
		n.Leaf.RemoveValue(value)
		return
	}

	currentKeywordTerm := keywordTerms[0]
	restKeywordTerms := keywordTerms[1:]

	switch currentKeywordTerm.Type {
	case "string":
		if len(currentKeywordTerm.Term) == 0 {
			n.RemoveString(restKeywordTerms, value)
			return
		}

		edgeCharacter := currentKeywordTerm.Term[:1]
		if currentKeywordTerm.CaseInsensitive {
			lowercase := strings.ToLower(edgeCharacter)
			uppercase := strings.ToUpper(edgeCharacter)
			if lowercase != uppercase {
				edgeCharacter = "x" + edgeCharacter
			}
		}
		rest := currentKeywordTerm.Term[1:]
		if node, ok := n.Edges[edgeCharacter]; ok {
			node.RemoveString(append(KeywordString{{Type: "string", Term: rest}}, restKeywordTerms...), value)
			if len(node.Edges) == 0 && !node.Leaf.HasValues() {
				delete(n.Edges, edgeCharacter)
			}
		}
	case "alphanumeric":
		fallthrough
	case "wordseparator":
		node, ok := n.Edges[currentKeywordTerm.Type]
		if ok {
			node.RemoveString(restKeywordTerms, value)
			if len(node.Edges) == 0 && !node.Leaf.HasValues() {
				delete(n.Edges, currentKeywordTerm.Type)
			}
		}
	case "alphanumeric*":
		fallthrough
	case "wordseparator*":
		node, ok := n.Edges[currentKeywordTerm.Type]
		if ok {
			node.RemoveString(restKeywordTerms, value)
			if len(node.Edges) == 1 && !node.Leaf.HasValues() {
				delete(n.Edges, currentKeywordTerm.Type)
			}
		}
	}
}

func (n *Node) Search(text string) []*Leaf {
	var asRune rune
	rest := text
	leaves := map[*Leaf]struct{}{}
	processingNodes := map[*Node]struct{}{n: {}}
	var newProcessingNodes map[*Node]struct{}
	processEdge := func(pNode *Node, edgeCharacter string) {
		if node, ok := pNode.Edges[edgeCharacter]; ok {
			if node.Leaf.HasValues() {
				leaves[node.Leaf] = struct{}{}
			}
			newProcessingNodes[node] = struct{}{}
		}
	}

	firstRun := true

	for rest != "" {
		edgeCharacter := rest[:1]
		rest = rest[1:]

		for _, c := range edgeCharacter {
			asRune = c
		}
		newProcessingNodes = map[*Node]struct{}{n: {}}
		if firstRun {
			firstRun = false
			if !isWordSeparator(asRune) {
				processEdge(n, "wordseparator")
				processEdge(n, "wordseparator*")
				processingNodes = newProcessingNodes
				newProcessingNodes = map[*Node]struct{}{n: {}}
			}
		}
		for pNode := range processingNodes {
			processEdge(pNode, edgeCharacter)
			lowercase := strings.ToLower(edgeCharacter)
			processEdge(pNode, "x"+lowercase)

			if isAlphanumeric(asRune) {
				processEdge(pNode, "alphanumeric")
				processEdge(pNode, "alphanumeric*")
			}
			if isWordSeparator(asRune) {
				processEdge(pNode, "wordseparator")
				processEdge(pNode, "wordseparator*")
			}
		}
		processingNodes = newProcessingNodes
	}

	if !isWordSeparator(asRune) {
		for pNode := range processingNodes {
			processEdge(pNode, "wordseparator")
			processEdge(pNode, "wordseparator*")
		}
	}

	return mapToSlice(leaves)
}
