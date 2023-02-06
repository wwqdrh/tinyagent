package srv

import (
	"fmt"

	"github.com/wwqdrh/tinyagent/pkg/srv/base"
	"github.com/wwqdrh/tinyagent/pkg/srv/redis"
)

var (
	SrvMap = map[string]func() base.ISrv{
		"redis": redis.NewRedisSrv,
	}

	srvList = []string{"redis", "mysql"}
)

func GetSrv(name string) (base.ISrv, error) {
	if v, ok := SrvMap[name]; !ok {
		return nil, fmt.Errorf("no this srv %s", name)
	} else {
		return v(), nil
	}
}

func SrvList() []string {
	return srvList
}
