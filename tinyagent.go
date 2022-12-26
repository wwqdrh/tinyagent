package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/tinyagent/api"
	"go.uber.org/zap"
)

var (
	Addr = flag.Int("addr", 8000, "服务监听端口")
)

func init() {
	flag.Parse()
}

func register(mux *http.ServeMux) {
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

func main() {
	mux := http.NewServeMux()
	register(mux)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *Addr),
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	logger.DefaultLogger.Infox("Starting httpserver at :%d\n", nil, *Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.DefaultLogger.Info("api exit...")
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.DefaultLogger.Error("server shutdown error", zap.Error(err))
	}
}
