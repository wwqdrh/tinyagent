package stack

type StackYamlConfig struct {
	File     string `yaml:"file"`
	External bool   `yaml:"external"`
}

func (s *StackYamlConfig) Active() bool {
	return false
}
