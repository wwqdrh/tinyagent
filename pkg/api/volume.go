package api

import (
	"net/http"

	"github.com/docker/docker/api/types/volume"
	"github.com/wwqdrh/tinyagent/agent/docker"
)

type VolumeListRes struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func VolumeList(w http.ResponseWriter, r *http.Request) {
	res, err := docker.VolumeList(volume.ListOptions{})
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

	err := docker.VolumeAdd(docker.VolumeSpec{
		Driver: req.Driver,
		Name:   req.Name,
	})
	if err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, "ok")
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

	if err := docker.VolumeDelete(req.Name, req.Force); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}
