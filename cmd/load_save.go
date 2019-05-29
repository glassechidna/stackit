package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

var DefaultOutputPath = "./stackit.packaged.yml"

func savePreparedOutput(prepared *stackit.PrepareOutput) error {
	jsonBody, err := ensureJson(prepared.TemplateBody)
	if err != nil {
		return errors.Wrap(err, "ensuring template is json")
	}

	prettyBuf := &bytes.Buffer{}
	err = json.Indent(prettyBuf, jsonBody, "", "  ")
	if err != nil {
		return errors.Wrap(err, "pretty-printing json of template body")
	}

	prepared.TemplateBody = prettyBuf.String()

	marshalled, err := yaml.Marshal(prepared)
	if err != nil {
		return errors.Wrap(err, "converting prepared output to yaml")
	}

	err = ioutil.WriteFile(DefaultOutputPath, marshalled, 0755)
	return errors.Wrapf(err, "writing prepared output to '%s'", DefaultOutputPath)
}

func loadPreparedOutput() (*stackit.PrepareOutput, error) {
	bytes, err := ioutil.ReadFile(DefaultOutputPath)
	if err != nil {
		return nil, errors.Wrapf(err, "reading prepared output from '%s'", DefaultOutputPath)
	}

	output := stackit.PrepareOutput{}
	err = yaml.Unmarshal(bytes, &output)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling prepared output yaml")
	}

	return &output, nil
}

func ensureJson(input string) ([]byte, error) {
	if strings.HasPrefix(input, "{") {
		return []byte(input), nil
	} else {
		// we assume it's yaml
		node := yaml.Node{}
		err := yaml.Unmarshal([]byte(input), &node)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling yaml into node")
		}

		return json.Marshal(jsonYamlNode(node))
	}
}

type jsonYamlNode yaml.Node

func (j jsonYamlNode) MarshalJSON() ([]byte, error) {
	switch j.Kind {
	case yaml.DocumentNode:
		return json.Marshal(jsonYamlNode(*j.Content[0]))
	case yaml.SequenceNode:
		arr := []jsonYamlNode{}
		for _, n := range j.Content {
			arr = append(arr, jsonYamlNode(*n))
		}
		return json.Marshal(arr)
	case yaml.MappingNode:
		kv := map[string]jsonYamlNode{}
		idx := 0
		for idx < len(j.Content) {
			keynode := j.Content[idx]
			valnode := j.Content[idx+1]
			kv[keynode.Value] = jsonYamlNode(*valnode)
			idx += 2
		}
		return json.Marshal(kv)
	case yaml.ScalarNode:
		// TODO what about numerics?
		return json.Marshal(j.Value)
	case yaml.AliasNode:
		panic("aliases not supported")
	}
	return nil, nil
}
