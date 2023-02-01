package srv

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/wwqdrh/gokit/logger"
)

var (
	SrvMap = map[string]func() ISrv{
		"redis": NewRedisSrv,
	}
)

type ISrv interface {
	SrvName() string
	SetOpt(name string, val interface{})
	Start()
	Stop()
}

type CurCmd struct {
	Base string
	Args []string
}

func GetSrv(name string) (ISrv, error) {
	if v, ok := SrvMap[name]; !ok {
		return nil, fmt.Errorf("no this srv %s", name)
	} else {
		return v(), nil
	}
}

func RunAndWait(cmd *exec.Cmd) (string, string, error) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	logger.DefaultLogger.Debugx("Task %s with args %+v", nil, cmd.Path, cmd.Args)
	err := cmd.Run()
	return outBuf.String(), errBuf.String(), err
}
