package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/tinyagent/agent/swarm"
	"github.com/wwqdrh/tinyagent/pkg/srv"
)

var (
	startOpt = struct {
		List bool `name:"list" alias:"l" desc:"列出所有支持的服务列表"`
	}{
		List: false,
	}

	stopOpt = struct {
		Name string `name:"name" required:"true"`
	}{}

	startMysqlOpt = struct {
		Name     string `name:"name"`
		Image    string `name:"image"`
		Port     int    `name:"port" desc:"example 3306"`
		Password string `name:"password"`
		Network  string `name:"network"`
	}{
		Name:     "mysql8",
		Image:    "bitnami/mysql:8.0",
		Password: "123456",
		Network:  "dev",
	}

	startRedisOpt = struct {
		Name     string `name:"name"`
		Image    string `name:"image"`
		Port     int    `name:"port"`
		Password string `name:"password" desc:"example:6379"`
		Network  string `name:"network"`
	}{
		Name:     "redis6",
		Image:    "bitnami/redis:6.2",
		Password: "123456",
		Network:  "dev",
	}

	startZipkinOpt = struct {
		Name     string `name:"name"`
		Image    string `name:"image"`
		Port     int    `name:"port"`
		Password string `name:"password" desc:"example:6379"`
		Network  string `name:"network"`
	}{
		Name:    "zipkin",
		Image:   "openzipkin/zipkin",
		Network: "dev",
	}

	startEtcdOpt = struct {
		Name         string `name:"name"`
		Image        string `name:"image"`
		ClientPort   int    `name:"clientport" desc:"example:2379,用于客户端链接"`
		PeerPort     int    `name:"peerport" desc:"example:2380,用于集群之间链接"`
		Password     string `name:"password"`
		Network      string `name:"network"`
		AdvertiseUrl string `name:"advertiseurl"`
	}{
		Name:         "etcd",
		Image:        "bitnami/etcd:3.5",
		AdvertiseUrl: "http://etcd:2379",
		Network:      "dev",
	}

	descOpt = struct {
		Name string
	}{}
)

func NewSrvStartCommand() *clitool.Command {
	cmd := &clitool.Command{
		Cmd: &cobra.Command{
			Use: "start",
			RunE: func(cmd *cobra.Command, args []string) error {
				if startOpt.List {
					fmt.Println("支持的服务列表如下: \nmysql\nredis")
					return nil
				} else {
					return cmd.Usage()
				}
			},
		},
		Values: &startOpt,
	}

	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "redis",
			RunE: func(cmd *cobra.Command, args []string) error {
				return (&srv.BitnamiRedisOpt{
					BaseSrvOpt: srv.BaseSrvOpt{
						Name:     startRedisOpt.Name,
						Image:    startRedisOpt.Image,
						Password: startRedisOpt.Password,
						Network:  startRedisOpt.Network,
						Ports:    map[int]int{startRedisOpt.Port: 6379},
					},
				}).Start()
			},
		},
		Values: &startRedisOpt,
	})

	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "zipkin",
			RunE: func(cmd *cobra.Command, args []string) error {
				return (&srv.ZipkinOpt{
					BaseSrvOpt: srv.BaseSrvOpt{
						Name:     startZipkinOpt.Name,
						Image:    startZipkinOpt.Image,
						Password: startZipkinOpt.Password,
						Network:  startZipkinOpt.Network,
						Ports:    map[int]int{startZipkinOpt.Port: startZipkinOpt.Port},
						Command:  []string{"start-zipkin"},
					},
				}).Start()
			},
		},
		Values: &startZipkinOpt,
	})

	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "mysql",
			RunE: func(cmd *cobra.Command, args []string) error {
				return (&srv.BitnamiMysqlOpt{
					Name:     startMysqlOpt.Name,
					Image:    startMysqlOpt.Image,
					Password: startMysqlOpt.Password,
					Network:  startMysqlOpt.Network,
					Ports:    map[int]int{startMysqlOpt.Port: 3306},
				}).Start()
			},
		},
		Values: &startMysqlOpt,
	})

	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "etcd",
			RunE: func(cmd *cobra.Command, args []string) error {
				return (&srv.BitnamiEtcdOpt{
					Name:                       startEtcdOpt.Name,
					Image:                      startEtcdOpt.Image,
					Password:                   startEtcdOpt.Password,
					Network:                    startEtcdOpt.Network,
					Ports:                      map[int]int{startEtcdOpt.ClientPort: 2379, startEtcdOpt.PeerPort: 2380},
					ETCD_ADVERTISE_CLIENT_URLS: startEtcdOpt.AdvertiseUrl,
				}).Start()
			},
		},
		Values: &startEtcdOpt,
	})

	return cmd
}

func NewSrvStopCommand() *clitool.Command {
	return &clitool.Command{
		Cmd: &cobra.Command{
			Use: "stop",
			RunE: func(cmd *cobra.Command, args []string) error {
				return swarm.ServiceRemove(stopOpt.Name)
			},
		},
		Values: &stopOpt,
	}
}

func NewSrvDescCommand() *clitool.Command {
	return &clitool.Command{
		Cmd: &cobra.Command{
			Use: "desc",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("do desc command")
				return nil
			},
		},
		Options: []clitool.OptionConfig{
			{
				Target:      "Name",
				Name:        "name",
				Description: "获取服务描述的名称",
				Required:    true,
			},
		},
		Values: &descOpt,
	}
}
