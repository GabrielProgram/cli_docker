package flags

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

type YamlFlag struct {
	raw []byte
}

func (y *YamlFlag) String() string {
	return fmt.Sprintf("YAML (%d bytes)", len(y.raw))
}

// TODO: Command.MarkFlagFilename()
func (y *YamlFlag) Set(v string) error {
	// Load request from file if it starts with '@' (like curl).
	if v[0] != '@' {
		y.raw = []byte(v)
		return nil
	}
	buf, err := os.ReadFile(v[1:])
	if err != nil {
		return fmt.Errorf("read %s: %w", v, err)
	}
	y.raw = buf
	return nil
}

func (y *YamlFlag) Unmarshal(v any) error {
	if y.raw == nil {
		return nil
	}
	return yaml.Unmarshal(y.raw, v)
}

func (y *YamlFlag) Type() string {
	return "YAML"
}
