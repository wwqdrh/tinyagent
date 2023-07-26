package srv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/tinyagent/pkg/swarm"
)

type BaseSrvOpt struct {
	Name     string
	Image    string
	Password string // if empty, equals ALLOW_EMPTY_PASSWORD=yes
	Network  string
	Ports    map[int]int
	Envs     []string
	Command  []string
	Configs  []swarm.ConfigOpt
}

func (o *BaseSrvOpt) SetOpt(name, value string) {
	switch name {
	case "name":
		o.Name = value
	case "image":
		o.Image = value
	case "password":
		o.Password = value
	case "network":
		o.Network = value
	case "port":
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return
		}
		p1, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			return
		}
		p2, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			return
		}
		o.Ports[int(p1)] = int(p2)
	case "command":
		o.Command = strings.Split(value, "\n")
	}
}

func (o *BaseSrvOpt) Start() error {
	res, err := swarm.ServiceCreate(
		swarm.ServiceOpt{
			Name:    o.Name,
			Image:   o.Image,
			Env:     o.Envs,
			Network: o.Network,
			Ports:   o.Ports,
			Command: o.Command,
			Configs: o.Configs,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(res.ID, res.Warnings)
	return nil
}
