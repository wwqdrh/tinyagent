package swarm

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type ServiceOpt struct {
	Name    string
	Image   string
	Env     []string
	Ports   map[int]int
	Network string
}

func (s *ServiceOpt) GetPorts() (res []swarm.PortConfig) {
	for hport, cport := range s.Ports {
		res = append(res, swarm.PortConfig{
			Name:          fmt.Sprintf("%s_%d", s.Name, hport),
			Protocol:      swarm.PortConfigProtocolTCP,
			TargetPort:    uint32(cport),
			PublishedPort: uint32(hport),
			PublishMode:   swarm.PortConfigPublishModeIngress,
		})
	}
	return
}

func ServiceExist(name string) (swarm.Service, []byte, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return swarm.Service{}, nil, err
	}
	defer cli.Close()

	return cli.ServiceInspectWithRaw(context.Background(), name, types.ServiceInspectOptions{})
}

func ServiceCreate(opt ServiceOpt) (types.ServiceCreateResponse, error) {
	if err := ImageExist(opt.Image); err != nil {
		if err != ErrImageNotExist {
			return types.ServiceCreateResponse{}, err
		}

		_, err := ImagePull(opt.Image, types.ImagePullOptions{})
		if err != nil {
			return types.ServiceCreateResponse{}, err
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return types.ServiceCreateResponse{}, err
	}
	defer cli.Close()

	return cli.ServiceCreate(context.Background(), swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: opt.Name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: opt.Image,
				Env:   opt.Env,
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{Target: opt.Network},
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Mode:  swarm.ResolutionModeVIP,
			Ports: opt.GetPorts(),
		},
	}, types.ServiceCreateOptions{})
}

func ServiceRemove(name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		return err
	}
	defer cli.Close()

	return cli.ServiceRemove(context.Background(), name)
}
