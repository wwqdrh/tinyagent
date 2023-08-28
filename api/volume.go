package api

import (
	"net/http"

	"github.com/docker/docker/api/types/volume"
	pkg_swarm "github.com/wwqdrh/tinyagent/pkg/swarm"
)

type VolumeListRes struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func VolumeList(w http.ResponseWriter, r *http.Request) {
	res, err := pkg_swarm.VolumeList(volume.ListOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}

	var result []VolumeListRes
	for _, item := range res.Volumes {
		result = append(result, VolumeListRes{
			Name:   item.Name,
			Driver: item.Driver,
		})
	}
	EchoJSON(w, ServerOK, result)
}

type VolumeCreateReq struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func VolumeCreate(w http.ResponseWriter, r *http.Request) {
	var req VolumeCreateReq
	if !ParseJSON(w, r, &req) {
		return
	}

	v, err := pkg_swarm.VolumeAdd(volume.CreateOptions{
		Driver: req.Driver,
		Name:   req.Name,
	})
	if err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, v)
	}
}

type VolumeRemoveReq struct {
	Name  string `json:"name"`
	Force bool   `json:"force"`
}

func VolumeRemove(w http.ResponseWriter, r *http.Request) {
	var req VolumeRemoveReq
	if !ParseJSON(w, r, &req) {
		return
	}

	if err := pkg_swarm.VolumeDelete(req.Name, req.Force); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}
