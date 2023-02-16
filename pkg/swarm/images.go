package swarm

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ImageExist(name string) error {
	images, err := ImageList(types.ImageListOptions{All: true})
	if err != nil {
		return err
	}

	for _, item := range images {
		for _, tag := range item.RepoTags {
			if tag == name {
				return nil
			}
		}
	}

	return ErrImageNotExist
}

func ImageList(opts types.ImageListOptions) ([]types.ImageSummary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	return cli.ImageList(context.Background(), opts)
}

func ImagePull(name string, opts types.ImagePullOptions) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return "", err
	}
	defer cli.Close()

	out, err := cli.ImagePull(context.Background(), name, opts)
	if err != nil {
		return "", err
	}

	buffer := bytes.Buffer{}
	buf := make([]byte, 100)
	for {
		n, err := out.Read(buf)
		if n > 0 {
			fmt.Println(string(buf[:n]))
			buffer.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
	}

	return buffer.String(), nil
}

func ImageDelete(name string, opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ImageRemove(context.Background(), name, opts)
}
