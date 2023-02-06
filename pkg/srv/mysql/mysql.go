package mysql

// -v./mysql/init:/docker-entrypoint-initdb.d/

import (
	"fmt"
	"strings"

	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/tinyagent/pkg/srv/base"
)

var (
	DefaultMysqlOpt = struct {
		Name     string
		Port     int
		Password string
		Mode     string
		Network  string
		Pull     bool
	}{}

	BasicOpt = struct {
		Image    string // mysql:8.0
		Name     string // 容器名字
		Port     int    // [port]:3306
		Password string // 管理员密码
		Pull     bool   // 是否拉取容器
	}{}

	SwarmOpt = struct {
		Image    string // mysql:8.0
		Name     string // 服务名称
		Password string
		Network  string // overlay network mode
		Pull     bool   // 是否拉取容器
	}{}
)

type MysqlSrv struct {
	mode string
}

func NewMysqlSrv() base.ISrv {
	return &MysqlSrv{
		mode: "basic",
	}
}

func (s MysqlSrv) SrvName() string {
	return "mysql"
}

func (s *MysqlSrv) SetOpt(name string, val interface{}) {

}

func (s MysqlSrv) Start() {
	// base.RunAndWait(exec.Command(
	// 	"docker", "pull"
	// ))
	if s.mode == "basic" {
		s.startBasic()
	} else if s.mode == "swarm" {
		s.startSwarm()
	} else {
		logger.DefaultLogger.Warn("unkonwn mode")
	}
}

func (s MysqlSrv) startBasic() {
	if BasicOpt.Pull {
		if err := base.RunAndWait2("docker", "pull", BasicOpt.Image); err != nil {
			return
		}
	}

	if err := base.RunAndWait2("docker", strings.Split(
		fmt.Sprintf(
			"run -d --name %s -p %d:3306 -e MYSQL_ROOT_PASSWORD=%s %s",
			BasicOpt.Name,
			BasicOpt.Port,
			BasicOpt.Password,
			BasicOpt.Image,
		), " ")...); err != nil {
		return
	}
}

func (s MysqlSrv) startSwarm() {
	if SwarmOpt.Pull {
		if err := base.RunAndWait2("docker", "pull", SwarmOpt.Image); err != nil {
			return
		}
	}

	if err := base.RunAndWait2("docker", strings.Split(
		fmt.Sprintf(
			"service create --name %s -e MYSQL_ROOT_PASSWORD=%s %s",
			SwarmOpt.Name,
			SwarmOpt.Password,
			SwarmOpt.Image,
		), " ")...); err != nil {
		return
	}
}

func (s MysqlSrv) Stop() {

}
