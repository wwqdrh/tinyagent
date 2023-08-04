package swarm

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func GetServiceIP(domain string) string {
	// Create a custom resolver with a dialer that connects to 127.0.0.11:53
	resolver := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.DialTimeout("udp", "127.0.0.11:53", 1*time.Second)
		},
	}

	// Use the resolver to lookup the host
	ips, err := resolver.LookupHost(context.TODO(), domain)
	if err != nil || len(ips) == 0 {
		return "127.0.0.1"
	}
	return ips[0]
}

func NetworkList(opts types.NetworkListOptions) ([]types.NetworkResource, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.NetworkList(context.Background(), opts)
}

func NetworkAdd(name, driver string) (types.NetworkCreateResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return types.NetworkCreateResponse{}, err
	}
	defer cli.Close()

	return cli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		IPAM: &network.IPAM{
			Driver: driver,
		},
	})
}

func NetworkRemove(name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
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
