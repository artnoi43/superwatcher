package utils

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadFileYAML[T any](filename string) (*T, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	var t T
	if err := yaml.Unmarshal(b, &t); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	return &t, nil
}
