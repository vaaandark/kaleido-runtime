package runtimeutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencontainers/runc/libcontainer"
)

func LoadLibcontainerState(root string, containerId string) (*libcontainer.State, error) {
	path := filepath.Join(root, containerId, "state.json")
	state := &libcontainer.State{}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, state)
	return state, err
}

func GetBundleFromState(state *libcontainer.State) string {
	for _, kv := range state.Config.Labels {
		split := strings.SplitN(kv, "=", 2)
		if split[0] == "bundle" {
			return split[1]
		}
	}
	return ""
}
