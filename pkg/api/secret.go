package api

import (
	"net/http"

	"github.com/docker/docker/api/types"
	pkg_swarm "github.com/wwqdrh/tinyagent/agent/swarm"
)

type SecretListRes struct {
	Name    string
	Content string
}

func SecretList(w http.ResponseWriter, r *http.Request) {
	res, err := pkg_swarm.SecretList(types.SecretListOptions{})
	if err != nil {
		EchoError(w, ServerError, err)
	} else {
		var result []SecretListRes
		for _, item := range res {
			result = append(result, SecretListRes{
				Name:    item.Spec.Name,
				Content: string(item.Spec.Data),
			})
		}
		EchoJSON(w, ServerOK, result)
	}
}

type SecretCreateReq struct {
	Name    string
	Content string
}

func SecretCreate(w http.ResponseWriter, r *http.Request) {
	var req SecretCreateReq
	if !ParseJSON(w, r, &req) {
		return
	}

	res, err := pkg_swarm.SecretAdd(req.Name, req.Content)
	if err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, res)
	}
}

type SecretRemoveReq struct {
	Name string `json:"name"`
}

func SecretRemove(w http.ResponseWriter, r *http.Request) {
	var req SecretRemoveReq
	if !ParseJSON(w, r, &req) {
		return
	}

	if err := pkg_swarm.SecretRemove(req.Name); err != nil {
		EchoError(w, ServerError, err)
	} else {
		EchoJSON(w, ServerOK, nil)
	}
}
