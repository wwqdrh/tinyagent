package swarm

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func VolumeList(filter volume.ListOptions) (volume.ListResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return volume.ListResponse{}, err
	}
	defer cli.Close()

	return cli.VolumeList(context.Background(), filter)
}

func VolumeAdd(options volume.CreateOptions) (volume.Volume, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return volume.Volume{}, err
	}
	defer cli.Close()

	return cli.VolumeCreate(context.Background(), options)
}

func VolumeDelete(name string, force bool) error {
	return withCli(func(cli *client.Client) error {
		return cli.VolumeRemove(context.Background(), name, force)
	})
}
