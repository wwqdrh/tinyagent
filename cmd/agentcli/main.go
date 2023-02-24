package main

import (
	"github.com/spf13/cobra"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/tinyagent/pkg/cli"
)

func main() {
	cmd := clitool.Command{
		Cmd: &cobra.Command{
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Usage()
			},
		},
	}
	cmd.Add(cli.NewSrvStartCommand())
	cmd.Add(cli.NewSrvStopCommand())
	cmd.Add(cli.NewSrvDescCommand())
	cmd.Run()
}
