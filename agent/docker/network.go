package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type NetworkSpec struct {
	Name    string   `json:"name"`
	Driver  string   `json:"driver"`
	Options []string `json:"options"`
	Subnet  string   `json:"subnet"`
	Gateway string   `json:"gateway"`
	IPRange string   `json:"ipRange"`
	Labels  []string `json:"labels"`
}

func NetworkList(opts types.NetworkListOptions) ([]types.NetworkResource, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.NetworkList(context.Background(), opts)
}

func NetworkAdd(spec NetworkSpec) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	var (
		ipam    network.IPAMConfig
		hasConf bool
	)
	if len(spec.Subnet) != 0 {
		ipam.Subnet = spec.Subnet
		hasConf = true
	}
	if len(spec.Gateway) != 0 {
		ipam.Gateway = spec.Gateway
		hasConf = true
	}
	if len(spec.IPRange) != 0 {
		ipam.IPRange = spec.IPRange
		hasConf = true
	}

	options := types.NetworkCreate{
		Driver:  spec.Driver,
		Options: stringsToMap(spec.Options),
		Labels:  stringsToMap(spec.Labels),
	}
	if hasConf {
		options.IPAM = &network.IPAM{Config: []network.IPAMConfig{ipam}}
	}
	if _, err := cli.NetworkCreate(context.TODO(), spec.Name, options); err != nil {
		return err
	}
	return nil
}

func NetworkRemove(name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return err
	}
	for _, network := range networks {
		if network.Name == name {
			return cli.NetworkRemove(context.Background(), network.ID)
		}
	}
	return errors.New("no this network")
}

func NetworkRemoveByIDs(ids ...string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	for _, id := range ids {
		if err := cli.NetworkRemove(context.TODO(), id); err != nil {
			if strings.Contains(err.Error(), "volume is in use") {
				return errors.Wrapf(err, "id: %s", id)
			}
			return err
		}
	}

	return nil
}
