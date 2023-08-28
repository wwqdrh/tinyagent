package api

import (
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/wwqdrh/tinyagent/agent/docker"
)

type NetworkListRes struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func NetworkList(w http.ResponseWriter, r *http.Request) {
	res, err := docker.NetworkList(types.NetworkListOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
	}
	var result []NetworkListRes
	for _, item := range res {
		result = append(result, NetworkListRes{
			Name:   item.Name,
			Driver: item.Driver,
		})
	}
	EchoJSON(w, ServerOK, result)
}

type NetworkCreateReq struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func NetworkCreate(w http.ResponseWriter, r *http.Request) {
	var req NetworkCreateReq
	if !ParseJSON(w, r, &req) {
		return
	}

	err := docker.NetworkAdd(docker.NetworkSpec{
		Name:   req.Name,
		Driver: req.Driver,
	})
	if err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, "ok")
	}
}

type NetworkRemoveReq struct {
	Name string `json:"name"`
}

func NetworkRemove(w http.ResponseWriter, r *http.Request) {
	var req NetworkRemoveReq
	if !ParseJSON(w, r, &req) {
		return
	}

	if err := docker.NetworkRemove(req.Name); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}
