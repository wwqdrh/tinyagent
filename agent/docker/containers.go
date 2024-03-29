package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	DefaultTimeout = 1 * time.Minute
)

func ContainerStart(name string, opts types.ContainerStartOptions) error {
	return withCli(func(cli *client.Client) error {
		return cli.ContainerStart(context.Background(), name, opts)
	})
}

// func ContainerRestart(name string) error {
// 	return withCli(func(cli *client.Client) error {
// 		return cli.ContainerRestart(context.Background(), name, &DefaultTimeout)
// 	})
// }

// func ContainerStop(name string) error {
// 	return withCli(func(cli *client.Client) error {
// 		return cli.ContainerStop(context.Background(), name, &DefaultTimeout)
// 	})
// }

func ContainerKill(name string) error {
	return withCli(func(cli *client.Client) error {
		return cli.ContainerKill(context.Background(), name, "KILL")
	})
}

func ContainerDelete(name string, opts types.ContainerRemoveOptions) error {
	return withCli(func(cli *client.Client) error {
		return cli.ContainerRemove(context.Background(), name, opts)
	})
}
