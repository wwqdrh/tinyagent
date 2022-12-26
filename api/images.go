package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	pkg_swarm "github.com/wwqdrh/tinyagent/pkg/swarm"
)

type ImageListRes struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

func ImageList(w http.ResponseWriter, r *http.Request) {
	images, err := pkg_swarm.ImageList(types.ImageListOptions{
		All: true,
	})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}

	var res []ImageListRes
	for _, item := range images {
		res = append(res, ImageListRes{
			ID:   item.ID,
			Tags: item.RepoTags,
		})
	}
	EchoJSON(w, ServerOK, res)
}

type ImagePullReq struct {
	Name string `json:"name"`
}

type ImagePullRes struct {
	Detail []string `json:"detail"`
}

func ImagePull(w http.ResponseWriter, r *http.Request) {
	var req ImagePullReq
	if !ParseJSON(w, r, &req) {
		return
	}

	out, err := pkg_swarm.ImagePull(req.Name, types.ImagePullOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}
	res, err := ioutil.ReadAll(out)
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}
	EchoJSON(w, ServerOK, ImagePullRes{
		Detail: strings.Split(string(res), "\n"),
	})
}

type ImageRemoveReq struct {
	Name string `json:"name"`
}

func ImageRemove(w http.ResponseWriter, r *http.Request) {
	var req ImageRemoveReq
	if !ParseJSON(w, r, &req) {
		return
	}

	resp, err := pkg_swarm.ImageDelete(req.Name, types.ImageRemoveOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}
	EchoJSON(w, ServerOK, resp)
}
