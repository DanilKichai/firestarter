package generator

import (
	"bytes"
	"firestarter/internal/pkg/batch"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v2"
)

func Load(file string, context interface{}) (*batch.Batch, error) {
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

	var u batch.Batch

	err = yaml.Unmarshal(buf.Bytes(), &u)
	if err != nil {
		return nil, fmt.Errorf("unmarshal batch: %w", err)
	}

	return &u, nil
}
