package swarm

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type ConfigOpt struct {
	Name   string
	Target string
}

func ConfigCreate(name string, data []byte) (types.ConfigCreateResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ConfigList(context.Background(), opts)
}

func ConfigGetByName(name string) (swarm.Config, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
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
