package discovery

import (
	"context"
	"sync"
	"time"

	"github.com/golang/glog"

	"go.etcd.io/etcd/client"
)

type NodesT map[string]string

func (w *Watcher) Len() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return len(w.Nodes)
}

func (w *Watcher) AllNodes() map[string]string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.Nodes
}

type Watcher struct {
	RootPath string
	KeysAPI  client.KeysAPI
	mutex    sync.Mutex
	Nodes    NodesT
}

func NewWatcher(endPoints []string, rootPath string) (*Watcher, error) {
	c, err := newClient(endPoints)
	if err != nil {
		return nil, err
	}

	keysAPI := client.NewKeysAPI(c)

	nodes := make(NodesT)

	w := &Watcher{
		RootPath: rootPath,
		KeysAPI:  keysAPI,
		Nodes:    nodes,
		mutex:    sync.Mutex{},
	}
	return w, nil
}

func (ns NodesT) addNode(key, value string) {
	ns[key] = value
}

func (ns NodesT) deleteNode(key string) {
	delete(ns, key)
}

func (w *Watcher) Watch() {
	api := w.KeysAPI
	for {
		resp, err := api.Get(context.Background(), w.RootPath, &client.GetOptions{Recursive: true})
		if err != nil {
			glog.Error(err)
			continue
		}

		if resp.Node != nil {
			w.mutex.Lock()
			w.Nodes = make(NodesT)
			for _, node := range resp.Node.Nodes {
				if node != nil {
					w.Nodes.addNode(node.Key, node.Value)
				}
			}
			w.mutex.Unlock()
		}

		time.Sleep(30 * time.Second)
	}

}
