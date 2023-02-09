package stack

type StackYamlNetwork struct {
	External bool `yaml:"external"`
}

func (s *StackYamlNetwork) Active() bool {
	return false
}
