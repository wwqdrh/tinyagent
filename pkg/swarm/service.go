package swarm

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/wwqdrh/gokit/logger"
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
		if hport == 0 {
			// 不对外暴露端口
			continue
		}
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

// 获取当前运行程序位置的service
// 需要指定是哪一个overlay环境下的
func CurrentService() (string, error) {
	curip, err := getClientIp("eth0")
	if err != nil {
		return "", err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(SupportedDockerAPIVersion))
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return "", err
	}
	defer cli.Close()

	networks, err := cli.NetworkList(context.TODO(), types.NetworkListOptions{})
	if err != nil {
		return "", err
	}

	for _, network := range networks {
		res, err := cli.NetworkInspect(context.TODO(), network.ID, types.NetworkInspectOptions{})
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
			return "", err
		}

		for _, srv := range res.Containers {
			if srv.IPv4Address == curip {
				return srv.Name, nil
			}
		}
	}

	return "", errors.New("not found")
}

// 寻找eth0这个网卡对应的ip
func getClientIp(device string) (string, error) {
	inter, err := net.InterfaceByName(device)
	if err != nil {
		return "", err
	}
	addrs, err := inter.Addrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.String(), nil
		}
	}

	return "", errors.New("can not find the client ip address")
}
