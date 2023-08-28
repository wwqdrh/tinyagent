package stack

import (
	"github.com/wwqdrh/gokit/logger"
	"gopkg.in/yaml.v3"
)

type StackYamlService struct {
	Image       string   `yaml:"image"`
	Command     string   `yaml:"command"`
	Networks    []string `yaml:"networks"`
	Environment []string `yaml:"environment"`
	Deploy      struct {
		Replicas int `yaml:"replicas" default:"1"`
	} `yaml:"deploy"`
	Configs []struct {
		Source string `yaml:"source"`
		Target string `yaml:"target"`
	} `yaml:"configs"`
	Volumes []string `yaml:"volumes"` // [name]:[target] 如果name只一个volume name那么type就是volume模式，否则就是bind模式
	Ports   []string `yaml:"ports"`
}

func NewServiceFromYamlString(content []byte) (res StackYamlService, err error) {
	err = yaml.Unmarshal(content, &res)
	if err != nil {
		logger.DefaultLogger.Warn(err.Error())
	}
	return
}

func (s *StackYamlService) Active() bool {
	return false
}
