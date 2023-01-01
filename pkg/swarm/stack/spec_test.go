package stack

import (
	"fmt"
	"testing"
)

func TestFromFile(t *testing.T) {
	v, err := NewStackYamlFromFile("./testdata/simple.yaml")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("+v", v)
	}
}
