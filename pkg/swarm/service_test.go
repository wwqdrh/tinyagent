package swarm

import (
	"fmt"
	"testing"
)

func TestServiceExist(t *testing.T) {
	s, data, err := ServiceExist("code-server")
	fmt.Println(s)
	fmt.Println(string(data))
	fmt.Println(err)

	s, data, err = ServiceExist("code-server1")
	fmt.Println(s)
	fmt.Println(string(data))
	fmt.Println(err)
}

func TestClientIP(t *testing.T) {
	ip, err := getClientIp("eth0")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(ip)
	}
}

func TestGetCurrentService(t *testing.T) {
	srvid, err := CurrentService()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(srvid)
	}
}
