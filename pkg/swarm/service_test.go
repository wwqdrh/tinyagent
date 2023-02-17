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
