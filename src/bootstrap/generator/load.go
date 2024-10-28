package generator

import (
	"bootstrap/unyaml"
	"bytes"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v2"
)

func Load(file string, context interface{}) (*unyaml.UnYAML, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	tmpl, err := template.New("bootstrap").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, context)
	if err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	var u unyaml.UnYAML

	err = yaml.Unmarshal(buf.Bytes(), &u)
	if err != nil {
		return nil, fmt.Errorf("unmarshal YAML: %w", err)
	}

	return &u, nil
}
