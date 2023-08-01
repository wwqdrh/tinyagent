package api

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	pkg_swarm "github.com/wwqdrh/tinyagent/pkg/swarm"
)

func ConfigCreate(w http.ResponseWriter, r *http.Request) {
	var content string
	if f, _, err := r.FormFile("content"); err != nil {
		EchoError(w, ClientParamInvalid, err)
		return
	} else {
		if data, err := ioutil.ReadAll(f); err != nil {
			EchoError(w, ClientParamInvalid, err)
			return
		} else {
			content = string(data)
		}
	}
	name := r.Form["name"]
	if len(name) != 1 {
		EchoError(w, ClientParamInvalid, errors.New("name参数为单数"))
		return
	}

	if resp, err := pkg_swarm.ConfigCreate(name[0], []byte(content)); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, H{"id": resp.ID})
	}
}

type ConfigListRes struct {
	Name    string `json:"name"`
	Version uint64 `json:"version"`
	Content string `json:"content"`
}

func ConfigList(w http.ResponseWriter, r *http.Request) {
	configs, err := pkg_swarm.ConfigList(types.ConfigListOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
		return
	}

	res := []ConfigListRes{}
	for _, item := range configs {
		res = append(res, ConfigListRes{
			Name:    item.Spec.Name,
			Version: item.Version.Index,
			Content: string(item.Spec.Data),
		})
	}
	EchoJSON(w, ServerOK, res)
}

func ConfigUpdate(w http.ResponseWriter, r *http.Request) {
	var content string
	if f, _, err := r.FormFile("content"); err != nil {
		EchoError(w, ClientParamInvalid, err)
		return
	} else {
		if data, err := ioutil.ReadAll(f); err != nil {
			EchoError(w, ClientParamInvalid, err)
			return
		} else {
			content = string(data)
		}
	}
	name := r.Form["name"]
	if len(name) != 1 {
		EchoError(w, ClientParamInvalid, errors.New("name参数为单数"))
		return
	}

	if err := pkg_swarm.ConfigUpdate(swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name: name[0],
		},
		Data: []byte(content),
	}); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}

type ConfigRemoveReq struct {
	Name string `json:"name"`
}

func ConfigRemove(w http.ResponseWriter, r *http.Request) {
	var req ConfigRemoveReq
	if !ParseJSON(w, r, &req) {
		return
	}

	if err := pkg_swarm.ConfigDelete(req.Name); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}
