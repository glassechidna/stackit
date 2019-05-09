package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var DefaultOutputPath = "./stackit.packaged.yml"

func savePreparedOutput(prepared *stackit.PrepareOutput) error {
	prettyBuf := &bytes.Buffer{}
	err := json.Indent(prettyBuf, []byte(prepared.TemplateBody), "", "  ")
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
