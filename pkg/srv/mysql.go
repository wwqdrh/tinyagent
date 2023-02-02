package srv

import (
	"fmt"
	"strings"
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
)

type MysqlSrv struct{}

func NewMysqlSrv() ISrv {
	return &MysqlSrv{}
}

func (s MysqlSrv) SrvName() string {
	return "mysql"
}

func (s *MysqlSrv) SetOpt(name string, val interface{}) {

}

func (s MysqlSrv) actions() (res []CurCmd) {
	if DefaultMysqlOpt.Pull {
		res = append(res, CurCmd{Base: "docker", Args: strings.Split("pull mysql:8.0", " ")})
	}

	if DefaultMysqlOpt.Mode == "basic" {
		res = append(res, CurCmd{Base: "docker", Args: strings.Split(fmt.Sprintf("run -d --name %s -v./mysql/init:/docker-entrypoint-initdb.d/ -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0", ""), " ")})
	}

	return
}

func (s MysqlSrv) Start() {

}

func (s MysqlSrv) Stop() {

}
