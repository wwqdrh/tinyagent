package srv

import "testing"

func TestStartRedisSrv(t *testing.T) {
	if err := (&BitnamiRedisOpt{
		BaseSrvOpt{
			Name:     "redis6",
			Image:    "bitnami/redis:6.2",
			Password: "123456",
			Network:  "dev",
			Ports:    map[int]int{6379: 6379},
		},
	}).Start(); err != nil {
		t.Error(err)
	}
}
