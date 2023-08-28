package api

import (
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/wwqdrh/tinyagent/agent/docker"
)

type ImageListRes struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

func ImageList(w http.ResponseWriter, r *http.Request) {
	images, err := docker.ImageList(types.ImageListOptions{
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

	out, err := docker.ImagePull(req.Name, types.ImagePullOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}
	EchoJSON(w, ServerOK, ImagePullRes{
		Detail: strings.Split(out, "\n"),
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

	resp, err := docker.ImageDelete(req.Name, types.ImageRemoveOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}
	EchoJSON(w, ServerOK, resp)
}
