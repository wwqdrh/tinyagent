package api

import (
	"net/http"

	"github.com/docker/docker/api/types"
	pkg_swarm "github.com/wwqdrh/tinyagent/pkg/swarm"
)

type NetworkListRes struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func NetworkList(w http.ResponseWriter, r *http.Request) {
	res, err := pkg_swarm.NetworkList(types.NetworkListOptions{})
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

	res, err := pkg_swarm.NetworkAdd(req.Name, req.Driver)
	if err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, res)
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

	if err := pkg_swarm.NetworkRemove(req.Name); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}
