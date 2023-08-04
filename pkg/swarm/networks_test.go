package swarm

import (
	"fmt"
	"testing"
)

func TestDomainResolve(t *testing.T) {
	if GetServiceIP("code-server-server") != "127.0.0.1" {
		t.Error("code-server-server解析失败")
	}

	fmt.Println(GetServiceIP("code-server"))
}
