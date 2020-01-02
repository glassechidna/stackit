package cfnyaml

import "gopkg.in/yaml.v3"

func Resolve(node *yaml.Node) {
	switch node.Kind {
	case yaml.DocumentNode:
		Resolve(node.Content[0])
	case yaml.SequenceNode:
		for _, n := range node.Content {
			Resolve(n)
		}
	case yaml.MappingNode:
		for _, n := range node.Content {
			Resolve(n)
		}

		var newcontent []*yaml.Node

		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			val := node.Content[i+1]
			if isMerge(key) {
				for _, sourceNode := range val.Alias.Content {
					newcontent = append(newcontent, sourceNode)
				}
				continue
			}

			if extantIdx := childExists(newcontent, key.Value); extantIdx >= 0 {
				newcontent[extantIdx] = key
				newcontent[extantIdx+1] = val
			} else {
				newcontent = append(newcontent, key)
				newcontent = append(newcontent, val)
			}
		}

		node.Content = newcontent

	case yaml.ScalarNode:
		break
	case yaml.AliasNode:
		Resolve(node.Alias)
	}

	node.Anchor = ""
}

func childExists(nodes []*yaml.Node, name string) int {
	for idx := 0; idx < len(nodes); idx += 2 {
		if nodes[idx].Value == name {
			return idx
		}
	}

	return -1
}

// copied from https://github.com/go-yaml/yaml/blob/f90ceb4f409096b60e2e9076b38b304b8246e5fa/decode.go#L813
func isMerge(n *yaml.Node) bool {
	return n.Kind == yaml.ScalarNode && n.Value == "<<" && (n.Tag == "" || n.Tag == "!" || n.Tag == "!!merge")
}
