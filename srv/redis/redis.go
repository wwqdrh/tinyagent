package redis

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/tinyagent/srv/base"
)

var (
	_ base.ISrv = &RedisSrv{}

	//go:embed redis.conf.template
	redisconf string

	DefaultRedisOpt = RedisSrvOpt{
		Name:     "redis6",
		Port:     6379,
		Password: "123456",
	}
)

type RedisSrv struct {
	opt *RedisSrvOpt
}

type RedisSrvOpt struct {
	Name     string
	Port     int
	Password string
	Mode     string
	Network  string
	Pull     bool
}

func NewRedisSrv() base.ISrv {
	var opt = DefaultRedisOpt
	return &RedisSrv{
		opt: &opt,
	}
}

func (s RedisSrv) SrvName() string {
	return "redis"
}

func (s *RedisSrv) SetOpt(name string, val interface{}) {
	switch name {
	case "name":
		s.opt.Name = val.(string)
	case "port":
		s.opt.Port = val.(int)
	case "password":
		s.opt.Password = val.(string)
	case "mode":
		s.opt.Mode = val.(string)
	case "network":
		s.opt.Network = val.(string)
	case "pull":
		s.opt.Pull = val.(bool)
	default:
		logger.DefaultLogger.Warnx("no this opt %s", nil, name)
	}
}

func (s RedisSrv) redisConf() (string, error) {
	type Conf struct {
		Password string
	}
	c := Conf{s.opt.Password}
	tmpl, err := template.New("conf").Parse(redisconf)
	if err != nil {
		return "", err
	}

	var out = bytes.NewBuffer([]byte{})
	err = tmpl.Execute(out, c)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile("./redis.conf", out.Bytes(), 0o777); err != nil {
		return "", err
	}
	return filepath.Abs("./redis.conf")
}

func (s RedisSrv) actions() (res []base.CurCmd) {
	conf, err := s.redisConf()
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return
	}

	if s.opt.Pull {
		res = append(res, base.CurCmd{Base: "docker", Args: strings.Split("pull redis:6-alpine", " ")})
	}

	if s.opt.Mode == "basic" {
		res = append(res, base.CurCmd{Base: "docker", Args: strings.Split(fmt.Sprintf(`run -d --name %s -v %s:/etc/redis/redis.conf -p %d:6379 redis:6-alpine redis-server /etc/redis/redis.conf`, s.opt.Name, conf, s.opt.Port), " ")})
	} else if s.opt.Mode == "swarm" {
		res = append(res, base.CurCmd{Base: "docker", Args: strings.Split(fmt.Sprintf(`service create --name %s --network %s --mount type=bind,src=%s,dst=/etc/redis/redis.conf redis:6-alpine redis-server /etc/redis/redis.conf`, s.opt.Name, s.opt.Network, conf), " ")})
	}

	return res
}

func (s RedisSrv) Start() {
	for _, item := range s.actions() {
		stdout, stderr, err := base.RunAndWait(exec.Command(item.Base, item.Args...))
		if err != nil || stderr != "" {
			logger.DefaultLogger.Error(err.Error())
			logger.DefaultLogger.Error(stderr)
			break
		}

		logger.DefaultLogger.Info(stdout)
	}
}

func (s RedisSrv) Stop() {
	actions := []base.CurCmd{
		{Base: "docker", Args: strings.Split(fmt.Sprintf("stop %s", s.opt.Name), " ")},
		{Base: "docker", Args: strings.Split(fmt.Sprintf("rm %s", s.opt.Name), " ")},
	}

	for _, item := range actions {
		stdout, stderr, err := base.RunAndWait(exec.Command(item.Base, item.Args...))
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
			break
		}
		if stderr != "" {
			logger.DefaultLogger.Error(err.Error())
			break
		}
		logger.DefaultLogger.Info(stdout)
	}
}
