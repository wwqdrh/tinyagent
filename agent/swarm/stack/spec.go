package stack

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/wwqdrh/gokit/basetool/datax"
	"github.com/wwqdrh/gokit/logger"
	"gopkg.in/yaml.v3"
)

type StackYaml struct {
	Version  string                      `yaml:"version"`
	Networks map[string]StackYamlNetwork `yaml:"networks"`
	Configs  map[string]StackYamlConfig  `yaml:"configs"`
	Volumes  map[string]StackYamlVolume  `yaml:"volumes"`
	Services map[string]StackYamlService `yaml:"services"`
}

func NewEmptyStackYaml() *StackYaml {
	return &StackYaml{
		Version:  "3.9",
		Networks: map[string]StackYamlNetwork{},
		Configs:  map[string]StackYamlConfig{},
		Volumes:  map[string]StackYamlVolume{},
		Services: map[string]StackYamlService{},
	}
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
	if err := res.LoadDefault(); err != nil {
		return StackYaml{}, err
	}
	return res, nil
}

func (c *StackYaml) Active() bool {
	for _, network := range c.Networks {
		network.Active()
	}

	for _, conf := range c.Configs {
		conf.Active()
	}

	for _, srv := range c.Services {
		srv.Active()
	}

	for _, volume := range c.Volumes {
		volume.Active()
	}

	return false
}

func (c *StackYaml) ExportYaml() (string, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		logger.DefaultLogger.Warn(err.Error())
		return "", err
	}
	return string(data), nil
}

func (c *StackYaml) AddService(name string, srv StackYamlService) {
	c.Services[name] = srv
}

func (c *StackYaml) AddServiceFromString(name string, data string) error {
	srv, err := NewServiceFromYamlString([]byte(data))
	if err != nil {
		return err
	}
	c.AddService(name, srv)
	return nil
}

func (c *StackYaml) AddNetwork(name string, external bool) {
	c.Networks[name] = StackYamlNetwork{
		External: external,
	}
}

func (c *StackYaml) AddVolume(name string, external bool) {
	c.Volumes[name] = StackYamlVolume{
		External: external,
	}
}

// TODO: 将basetool.datax里面的LoadDefault提供嵌套能力
func (c *StackYaml) LoadDefault() error {
	for name := range c.Services {
		srv := c.Services[name]
		if err := datax.LoadDefault(&srv); err != nil {
			return err
		}
		c.Services[name] = srv
	}
	return nil
}

func (c *StackYaml) GetConfig(name string) (StackYamlConfig, error) {
	for confname, conf := range c.Configs {
		if confname == name {
			return conf, nil
		}
	}
	return StackYamlConfig{}, ErrConfigNotExist
}

func (c *StackYaml) UpdateConfig(name string, conf StackYamlConfig) error {
	for confname := range c.Configs {
		if confname == name {
			c.Configs[name] = conf
			return nil
		}
	}
	return ErrConfigNotExist
}

func (c *StackYaml) GetServiceByConfig(name string) []string {
	res := []string{}
	for srvname, item := range c.Services {
		for _, conf := range item.Configs {
			if conf.Source == name {
				res = append(res, srvname)
			}
		}
	}
	return res
}

func (c *StackYaml) GetAndUpdateServiceByConfig(name string, newconf string) error {
	conf, err := c.GetConfig(name)
	if err != nil {
		return err
	}
	c.Configs[newconf] = conf
	delete(c.Configs, name)

	for srvname, item := range c.Services {
		for _, conf := range item.Configs {
			if conf.Source == name {
				if err := c.UpdateServiceConfig(srvname, name, newconf); err != nil {
					logger.DefaultLogger.Warn(err.Error())
				}
			}
		}
	}
	return nil
}

func (c *StackYaml) UpdateServiceConfig(srvname, oldconf, newconf string) error {
	srv, ok := c.Services[srvname]
	if !ok {
		return fmt.Errorf("%s不存在", srvname)
	}

	for i, conf := range srv.Configs {
		if conf.Source == oldconf {
			conf.Source = newconf
			srv.Configs[i] = conf
		}
	}
	return nil
}

func (c *StackYaml) GetServiceByVolume(name string) []string {
	res := []string{}
	for srvname, item := range c.Services {
		for _, volm := range item.Volumes {
			if name == strings.Split(volm, ":")[0] {
				res = append(res, srvname)
			}
		}
	}
	return res
}

func (c *StackYaml) GetAndUpdateServiceByVolume(name string, newvolume string) {
	for srvname, item := range c.Services {
		for _, volm := range item.Volumes {
			if name == strings.Split(volm, ":")[0] {
				if err := c.UpdateServiceVolume(srvname, name, newvolume); err != nil {
					logger.DefaultLogger.Warn(err.Error())
				}
			}
		}
	}
}

func (c *StackYaml) UpdateServiceVolume(srvname, oldvolume, newvolume string) error {
	srv, ok := c.Services[srvname]
	if !ok {
		return fmt.Errorf("%s不存在", srvname)
	}

	for i, volume := range srv.Volumes {
		t := strings.Split(volume, ":")
		if oldvolume == t[0] {
			srv.Volumes[i] = fmt.Sprintf("%s:%s", newvolume, t[1])
		}
	}
	return nil
}
