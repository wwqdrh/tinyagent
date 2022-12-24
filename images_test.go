package tinyagent

import (
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestImageList(t *testing.T) {
	res, err := ImageList(types.ImageListOptions{All: true})
	if err != nil {
		t.Error(err)
	}
	for _, item := range res {
		fmt.Printf("%s %v\n", item.ID, item.RepoTags)
	}
}
