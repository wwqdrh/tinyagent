package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/tinyagent/pkg/swarm"
	redis "github.com/wwqdrh/tinyagent/srv"
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
		Port     int    `name:"port"`
		Password string `name:"password"`
		Network  string `name:"network"`
	}{
		Name:     "mysql8",
		Image:    "bitnami/mysql:8.0",
		Port:     3306,
		Password: "123456",
	}

	startRedisOpt = struct {
		Name     string `name:"name"`
		Image    string `name:"image"`
		Port     int    `name:"port"`
		Password string `name:"password"`
		Network  string `name:"network"`
	}{
		Name:     "redis6",
		Image:    "bitnami/redis:6.2",
		Port:     6379,
		Password: "123456",
	}

	descOpt = struct {
		Name string
	}{}
)

func startCommand() *clitool.Command {
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
				return (&redis.BitnamiRedisOpt{
					Name:     startRedisOpt.Name,
					Image:    startRedisOpt.Image,
					Password: startRedisOpt.Password,
					Network:  startRedisOpt.Network,
					Ports:    map[int]int{startRedisOpt.Port: 6379},
				}).Start()
			},
		},
		Values: &startRedisOpt,
	})

	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "mysql",
			RunE: func(cmd *cobra.Command, args []string) error {
				return (&redis.BitnamiMysqlOpt{
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

	return cmd
}

func stopCommand() *clitool.Command {
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

func descCommand() *clitool.Command {
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

func main() {
	cmd := clitool.Command{
		Cmd: &cobra.Command{
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Usage()
			},
		},
	}
	cmd.Add(startCommand())
	cmd.Add(stopCommand())
	cmd.Add(descCommand())
	cmd.Run()
}
