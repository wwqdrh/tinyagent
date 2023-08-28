package swarm

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type ConfigOpt struct {
	Name   string
	Target string
}

type ConfigEntryFile struct {
	Source string
	Name   string
	Target string
}

func (c ConfigEntryFile) Create() error {
	data, err := os.ReadFile(c.Source)
	if err != nil {
		return errors.Wrapf(err, "读取文件失败")
	}
	_, err = ConfigCreate(c.Name, data)
	if err != nil {
		return errors.Wrapf(err, "创建配置文件失败")
	}
	return nil
}

type ConfigEntrySource struct {
	Data   []byte
	Name   string
	Target string
}

func (c ConfigEntrySource) Create() error {
	_, err := ConfigCreate(c.Name, c.Data)
	if err != nil {
		return errors.Wrapf(err, "创建配置文件失败")
	}
	return nil
}

func ConfigCreate(name string, data []byte) (types.ConfigCreateResponse, error) {
	if ok, err := ConfigExist(name); err != nil || !ok {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return types.ConfigCreateResponse{}, err
		}
		defer cli.Close()
		return cli.ConfigCreate(context.Background(), swarm.ConfigSpec{
			Annotations: swarm.Annotations{
				Name: name,
			},
			Data: data,
		})
	} else {
		return types.ConfigCreateResponse{}, errors.Wrapf(err, "配置已经存在")
	}
}

func ConfigExist(name string) (bool, error) {
	res, err := ConfigList(types.ConfigListOptions{})
	if err != nil {
		return false, err
	}
	for _, item := range res {
		if item.Spec.Name == name {
			return true, nil
		}
	}
	return false, errors.New("not exist this config")
}

func ConfigList(opts types.ConfigListOptions) ([]swarm.Config, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ConfigList(context.Background(), opts)
}

func ConfigGetByName(name string) (swarm.Config, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return swarm.Config{}, err
	}
	defer cli.Close()

	configs, err := cli.ConfigList(context.Background(), types.ConfigListOptions{})
	if err != nil {
		return swarm.Config{}, err
	}
	for _, item := range configs {
		if item.Spec.Name == name {
			return item, nil
		}
	}
	return swarm.Config{}, errors.New("config not exist")
}

func ConfigUpdate(spec swarm.ConfigSpec) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	conf, err := ConfigGetByName(spec.Name)
	if err != nil {
		return err
	}

	version := conf.Version
	version.Index += 1

	return cli.ConfigUpdate(context.Background(), conf.ID, version, spec)
}

func ConfigDelete(name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	conf, err := ConfigGetByName(name)
	if err != nil {
		return err
	}

	return cli.ConfigRemove(context.Background(), conf.ID)
}
