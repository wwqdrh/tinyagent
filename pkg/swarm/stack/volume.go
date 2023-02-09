package stack

type StackYamlVolume struct {
	External bool `yaml:"external"`
}

func (s *StackYamlVolume) Active() bool {
	return false
}
