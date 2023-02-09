package srv

import "testing"

func TestStartRedis(t *testing.T) {
	srv, err := GetSrv("redis")
	if err != nil {
		t.Error(err)
		return
	}

	defer srv.Stop()
	srv.Start()
}
