package swarm

import (
	"fmt"
	"testing"
)

func TestIsSwarm(t *testing.T) {
	ok, err := IsSwarm()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(ok)
	}
}
