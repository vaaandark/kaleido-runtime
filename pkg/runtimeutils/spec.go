package runtimeutils

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func LoadRuntimeSpecByBundle(bundle string) (*specs.Spec, error) {
	path := filepath.Join(bundle, "config.json")
	return LoadRuntimeSpec(path)
}

func LoadRuntimeSpec(path string) (*specs.Spec, error) {
	spec := &specs.Spec{}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, spec)
	return spec, err
}
