package base

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/wwqdrh/gokit/logger"
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

func RunAndWait(cmd *exec.Cmd) (string, string, error) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	logger.DefaultLogger.Debugx("Task %s with args %+v", nil, cmd.Path, cmd.Args)
	err := cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func RunAndWait2(command string, args ...string) error {
	outstr, errstr, err := RunAndWait(exec.Command(command, args...))
	if err != nil || errstr != "" {
		e := fmt.Errorf("%w: %s", err, errstr)
		logger.DefaultLogger.Warn(e.Error())
		return e
	}
	logger.DefaultLogger.Info(outstr)
	return nil
}
