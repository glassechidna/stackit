package cfnyaml

import (
	"bytes"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type CfnYaml struct {
	yaml.Node
}

func (c *CfnYaml) MarshalYAML() (interface{}, error) {
	return &c.Node, nil
}

func ParseFile(path string) (*CfnYaml, error) {
	c := &CfnYaml{}

	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}

	err = yaml.Unmarshal(byt, &c.Node)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling yaml")
	}

	return c, nil
}

func (c *CfnYaml) String() string {
	buf := &bytes.Buffer{}
	w := yaml.NewEncoder(buf)
	w.SetIndent(2)
	_ = w.Encode(&c.Node)
	return buf.String()
}

type PackageableNode struct {
	Name    string
	Value   string
	Replace func(bucket, key, versionId string)
	path    *yaml.Node
}

func (c *CfnYaml) PackageableNodes() ([]PackageableNode, error) {
	var nodes []PackageableNode

	resources := valueForKey(&c.Node, "Resources")
	if resources == nil {
		return nil, errors.New("no top-level key named `Resources` found in template")
	}

	idx := 0
	for idx < len(resources.Content) {
		nameNode := resources.Content[idx]
		name := nameNode.Value

		valueNode := resources.Content[idx+1]
		resTypeNode := valueForKey(valueNode, "Type")
		if resTypeNode == nil {
			return nil, errors.Errorf("resource `%s` has no `Type`", name)
		}
		resType := resTypeNode.Value

		if def := packageableDefinition(resType); def != nil {
			propNode := valueForKey(valueNode, def.Path...)
			if propNode != nil {
				nodes = append(nodes, PackageableNode{
					Name:  name,
					Value: propNode.Value,
					Replace: func(bucket, key, versionId string) {
						newNode := def.Rewritten(bucket, key, versionId)
						*propNode = *newNode
					},
				})
			}
		}

		idx += 2
	}

	return nodes, nil
}

func packageableDefinition(typ string) *packageablePropertyDefinition {
	for _, def := range packageablePropertyDefinitions {
		if def.ResourceType == typ {
			return &def
		}
	}

	return nil
}

func valueForKey(n *yaml.Node, key ...string) *yaml.Node {
	if n.Kind == yaml.DocumentNode {
		return valueForKey(n.Content[0], key...)
	} else if n.Kind != yaml.MappingNode {
		// TODO: panic?
		return nil
	}

	keyFound := false
	for _, c := range n.Content {
		if c.Value == key[0] {
			keyFound = true
		} else if keyFound {
			if len(key) == 1 {
				return c
			} else {
				return valueForKey(c, key[1:]...)
			}
		}
	}

	return nil
}
