package swarm

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func VolumeList(filter filters.Args) (volume.VolumeListOKBody, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return volume.VolumeListOKBody{}, err
	}
	defer cli.Close()

	return cli.VolumeList(context.Background(), filter)
}

func VolumeAdd(options volume.VolumeCreateBody) (types.Volume, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return types.Volume{}, err
	}
	defer cli.Close()

	return cli.VolumeCreate(context.Background(), options)
}

func VolumeDelete(name string, force bool) error {
	return withCli(func(cli *client.Client) error {
		return cli.VolumeRemove(context.Background(), name, force)
	})
}
