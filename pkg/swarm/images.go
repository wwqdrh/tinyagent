package swarm

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ImageList(opts types.ImageListOptions) ([]types.ImageSummary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	return cli.ImageList(context.Background(), opts)
}

func ImagePull(name string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ImagePull(context.Background(), name, opts)
}

func ImageDelete(name string, opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ImageRemove(context.Background(), name, opts)
}
