package stack

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type StackYaml struct {
	Version  string                      `yaml:"version"`
	Networks map[string]StackYamlNetwork `yaml:"networks"`
	Configs  map[string]StackYamlConfig  `yaml:"configs"`
	Volumes  map[string]StackYamlVolume  `yaml:"volumes"`
}

type StackYamlNetwork struct {
	External bool `yaml:"external"`
}

type StackYamlConfig struct {
	File     string `yaml:"file"`
	External bool   `yaml:"external"`
}

type StackYamlVolume struct {
	External bool `yaml:"external"`
}

type StackYamlService struct {
	Image       string   `yaml:"image"`
	Command     []string `yaml:"command"`
	Networks    []string `yaml:"networks"`
	Environment []string `yaml:"environment"`
	Deploy      struct {
		Replicas int `yaml:"replicas"`
	} `yaml:"deploy"`
	Configs []struct {
		Source string `yaml:"source"`
		Target string `yaml:"target"`
	} `yaml:"configs"`
	Ports []string `yaml:"ports"`
}

func NewStackYamlFromFile(f string) (StackYaml, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return StackYaml{}, err
	}
	var res StackYaml
	if err = yaml.Unmarshal(data, &res); err != nil {
		return StackYaml{}, err
	}
	return res, nil
}
