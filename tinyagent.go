package tinyagent

import (
	"flag"
	"net/http"

	"github.com/wwqdrh/tinyagent/pkg/api"
)

var (
	Addr = flag.Int("addr", 8000, "服务监听端口")
)

type HttpEngine interface {
	HandleFunc(string, http.HandlerFunc)
}

func RegisterAPI(mux HttpEngine) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/swarm/config/list", api.ConfigList)
	mux.HandleFunc("/swarm/config/create", api.ConfigCreate)
	mux.HandleFunc("/swarm/config/update", api.ConfigUpdate)
	mux.HandleFunc("/swarm/config/remove", api.ConfigRemove)
	mux.HandleFunc("/swarm/image/list", api.ImageList)
	mux.HandleFunc("/swarm/image/pull", api.ImagePull)
	mux.HandleFunc("/swarm/image/remove", api.ImageRemove)
	mux.HandleFunc("/swarm/volume/list", api.VolumeList)
	mux.HandleFunc("/swarm/volume/add", api.VolumeCreate)
	mux.HandleFunc("/swarm/volume/remove", api.VolumeRemove)
	mux.HandleFunc("/swarm/secret/list", api.SecretList)
	mux.HandleFunc("/swarm/secret/create", api.SecretCreate)
	mux.HandleFunc("/swarm/secret/remove", api.SecretRemove)
	mux.HandleFunc("/swarm/network/list", api.NetworkList)
	mux.HandleFunc("/swarm/network/create", api.NetworkCreate)
	mux.HandleFunc("/swarm/network/remove", api.NetworkRemove)
}
