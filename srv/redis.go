package srv

import (
	"fmt"
	"strconv"

	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/tinyagent/pkg/swarm"
)

// env
// REDIS_PASSWORD=password123
// ALLOW_EMPTY_PASSWORD=yes
//
type BitnamiRedisOpt struct {
	Name     string
	Image    string
	Password string // if empty, equals ALLOW_EMPTY_PASSWORD=yes
	Network  string
	Ports    map[int]int
}

func (o *BitnamiRedisOpt) SetOpt(name, value string) {
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
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
		} else {
			o.Ports[int(v)] = 6379
		}
	}
}

func (o *BitnamiRedisOpt) envs() (res []string) {
	if o.Password == "" {
		res = append(res, "ALLOW_EMPTY_PASSWORD=yes")
	} else {
		res = append(res, "REDIS_PASSWORD="+o.Password)
	}
	return res
}

func (o *BitnamiRedisOpt) Start() error {
	if o.Image == "" {
		o.Image = "bitnami/redis:6.2"
	}

	res, err := swarm.ServiceCreate(
		swarm.ServiceOpt{
			Name:    o.Name,
			Image:   o.Image,
			Env:     o.envs(),
			Network: o.Network,
			Ports:   o.Ports,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(res.ID, res.Warnings)
	return nil
}
