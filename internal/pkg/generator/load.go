package generator

import (
	"bytes"
	"firestarter/internal/pkg/batch"
	"firestarter/internal/pkg/generator/helpers"
	"fmt"
	"html/template"
	"path"

	"gopkg.in/yaml.v2"
)

func Load(file string, context interface{}, show bool) (*batch.Batch, error) {

	funcMap := template.FuncMap{
		"path2uri": helpers.Path2URI,
	}

	name := path.Base(file)
	tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(file)
	if err != nil {
		return nil, fmt.Errorf("construct template: %w", err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, context)
	if err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	if show {
		fmt.Print(buf.String())
	}

	var u batch.Batch

	err = yaml.Unmarshal(buf.Bytes(), &u)
	if err != nil {
		return nil, fmt.Errorf("unmarshal batch: %w", err)
	}

	return &u, nil
}
