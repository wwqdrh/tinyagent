package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/tinyagent/pkg/srv"
)

var (
	startOpt = struct {
		Name string
	}{}

	startRedisOpt = struct {
		Name     string
		Port     int
		Password string
	}{}

	descOpt = struct {
		Name string
	}{}
)

func startCommand() *clitool.Command {
	cmd := &clitool.Command{
		Cmd: &cobra.Command{
			Use: "start",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Usage()
			},
		},
		Options: []clitool.OptionConfig{
			{
				Target:      "Name",
				Name:        "name",
				Description: "启动的服务名称",
				Required:    true,
			},
		},
		Values: &startOpt,
	}

	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "redis",
			RunE: func(cmd *cobra.Command, args []string) error {
				factory, err := srv.GetSrv("redis")
				if err != nil {
					return err
				}

				factory.SetOpt("name", startRedisOpt.Name)
				factory.SetOpt("port", startRedisOpt.Port)
				factory.SetOpt("password", startRedisOpt.Password)
				factory.Start()
				return nil
			},
		},
		Options: []clitool.OptionConfig{
			{
				Target:       "Name",
				Name:         "name",
				Description:  "启动的服务名称",
				DefaultValue: "redis6",
			},
			{
				Target:       "Port",
				Name:         "port",
				Description:  "启动暴露的端口",
				DefaultValue: 6379,
			},
			{
				Target:       "Password",
				Name:         "password",
				Description:  "redis密码",
				DefaultValue: "123456",
			},
		},
		Values: &startRedisOpt,
	})

	return cmd
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

	cmd.Add(&clitool.Command{
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
	})

	cmd.Run()
}
