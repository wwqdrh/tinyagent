package swarm

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func SecretList(opts types.SecretListOptions) ([]swarm.Secret, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.SecretList(context.Background(), opts)
}

func SecretAdd(name, content string) (types.SecretCreateResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.SecretCreateResponse{}, err
	}
	defer cli.Close()

	return cli.SecretCreate(context.Background(), swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		Data: []byte(content),
	})
}

func SecretRemove(name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	secrets, err := cli.SecretList(context.Background(), types.SecretListOptions{})
	if err != nil {
		return err
	}
	for _, secret := range secrets {
		if secret.Spec.Name == name {
			return cli.SecretRemove(context.Background(), secret.ID)
		}
	}
	return errors.New("no this secret")
}
