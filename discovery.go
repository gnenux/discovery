package discovery

import (
	"time"

	"go.etcd.io/etcd/client"
)

var (
	DefaultWatcher *Watcher
)

func newClient(endPoints []string) (client.Client, error) {
	cfg := client.Config{
		Endpoints:               endPoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 60 * time.Second,
	}

	return client.New(cfg)
}

func InitDefaultWatcher(endPoints []string, rootPath string) error {
	w, err := NewWatcher(endPoints, rootPath)
	if err != nil {
		return err
	}

	DefaultWatcher = w
	return nil
}
