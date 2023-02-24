package srv

import (
	"fmt"

	"github.com/wwqdrh/tinyagent/pkg/swarm"
)

type BitnamiEtcdOpt struct {
	Name                       string
	Image                      string
	Password                   string
	Network                    string
	Ports                      map[int]int
	ETCD_ADVERTISE_CLIENT_URLS string // 广播给其他成员的地址
}

func (o *BitnamiEtcdOpt) SetOpt(name, value string) {
	switch name {
	case "name":
		o.Name = value
	case "image":
		o.Image = value
	case "network":
		o.Network = value
	case "advertiseurl":
		o.ETCD_ADVERTISE_CLIENT_URLS = value
	}
}

func (o *BitnamiEtcdOpt) envs() (res []string) {
	if o.Password == "" {
		res = append(res, "ALLOW_NONE_AUTHENTICATION=yes")
	} else {
		res = append(res, "ETCD_ROOT_PASSWORD="+o.Password)
	}

	res = append(res,
		"ETCD_ADVERTISE_CLIENT_URLS="+o.ETCD_ADVERTISE_CLIENT_URLS,
	)
	return
}

func (o *BitnamiEtcdOpt) Start() error {
	if o.Image == "" {
		o.Image = "bitnami/etcd:3.5"
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
