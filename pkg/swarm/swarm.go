package swarm

import (
	"context"
	"errors"

	"github.com/docker/docker/client"
	"github.com/wwqdrh/gokit/logger"
)

var (
	ErrImageNotExist = errors.New("no this image")
)

func IsSwarm() (bool, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return false, err
	}
	defer cli.Close()

	_, err = cli.SwarmInspect(context.Background())
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return false, err
	}

	return true, nil
}
