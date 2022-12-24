package tinyagent

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wwqdrh/gokit/logger"
	"go.uber.org/zap"
)

const (
	serviceNameLabel = "com.docker.swarm.service.name"
)

type InfoService struct{}

func NewInfoService() *InfoService {
	return &InfoService{}
}

// GetRuntimeConfigurationFromDockerEngine retrieves information from a Docker environment
// and returns a map of labels.
func (service *InfoService) GetRuntimeConfigurationFromDockerEngine() (*RuntimeConfiguration, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	dockerInfo, err := cli.Info(context.Background())
	if err != nil {
		return nil, err
	}

	runtimeConfiguration := &RuntimeConfiguration{
		NodeName:            dockerInfo.Name,
		DockerConfiguration: DockerRuntimeConfiguration{},
	}

	if dockerInfo.Swarm.NodeID == "" {
		getStandaloneConfiguration(runtimeConfiguration)
	} else {
		err := getSwarmConfiguration(runtimeConfiguration, dockerInfo, cli)
		if err != nil {
			return nil, err
		}

	}

	return runtimeConfiguration, nil
}

// GetContainerIpFromDockerEngine is used to retrieve the IP address of the container through Docker.
// It will inspect the container to retrieve the networks associated to the container and returns the IP associated
// to the first network found that is not an ingress network. If the ignoreNonSwarmNetworks parameter is specified,
// it will also ignore non Swarm scoped networks.
func (service *InfoService) GetContainerIpFromDockerEngine(containerName string, ignoreNonSwarmNetworks bool) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return "", err
	}
	defer cli.Close()

	containerInspect, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return "", err
	}

	if len(containerInspect.NetworkSettings.Networks) > 1 {
		logger.DefaultLogger.Warn("agent container running in more than a single Docker network. This might cause communication issues",
			zap.Int("network_count", len(containerInspect.NetworkSettings.Networks)))
	}

	for networkName, network := range containerInspect.NetworkSettings.Networks {
		networkInspect, err := cli.NetworkInspect(context.Background(), network.NetworkID, types.NetworkInspectOptions{})
		if err != nil {
			return "", err
		}

		if networkInspect.Ingress || (ignoreNonSwarmNetworks && networkInspect.Scope != "swarm") {
			logger.DefaultLogger.Debug("skipping invalid container network",
				zap.String("network_name", networkInspect.Name),
				zap.String("scope", networkInspect.Scope),
				zap.Bool("ingress", networkInspect.Ingress),
			)
			continue
		}

		if network.IPAddress != "" {
			logger.DefaultLogger.Debug("retrieving IP address from container network",
				zap.String("ip_address", network.IPAddress),
				zap.String("network_name", networkName),
			)

			return network.IPAddress, nil
		}
	}

	return "", errors.New("unable to retrieve the address on which the agent can advertise. Check your network settings")
}

// GetServiceNameFromDockerEngine is used to return the name of the Swarm service the agent is part of.
// The service name is retrieved through container labels.
func (service *InfoService) GetServiceNameFromDockerEngine(containerName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return "", err
	}
	defer cli.Close()

	containerInspect, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return "", err
	}

	return containerInspect.Config.Labels[serviceNameLabel], nil
}

func getStandaloneConfiguration(config *RuntimeConfiguration) {
	config.DockerConfiguration.EngineStatus = EngineStatusStandalone
}

func getSwarmConfiguration(config *RuntimeConfiguration, dockerInfo types.Info, cli *client.Client) error {
	config.DockerConfiguration.EngineStatus = EngineStatusSwarm
	config.DockerConfiguration.NodeRole = NodeRoleWorker

	if dockerInfo.Swarm.ControlAvailable {
		config.DockerConfiguration.NodeRole = NodeRoleManager

		node, _, err := cli.NodeInspectWithRaw(context.Background(), dockerInfo.Swarm.NodeID)
		if err != nil {
			return err
		}

		if node.ManagerStatus.Leader {
			config.DockerConfiguration.Leader = true
		}
	}

	return nil
}

func withCli(callback func(cli *client.Client) error) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return err
	}
	defer cli.Close()

	return callback(cli)
}
