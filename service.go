package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gnenux/xlog"
	"go.etcd.io/etcd/client"
)

type Service struct {
	ProcessId int        // 进程ID
	RootPath  string     //根目录
	Info      ServerInfo // 服务端信息
	KeysAPI   client.KeysAPI
}

type ServerInfo struct {
	Id   int    `json:"id"`   // 服务器ID
	Name string `json:"name"` //服务器name
	IP   string `json:"ip"`   // 对外连接服务的IP
	Port int    `json:"port"` // 对外服务端口
}

func (s *ServerInfo) String() string {
	return fmt.Sprintf("%s:%s:%d", s.Name, s.IP, s.Port)
}

func RegisterService(endPoints []string, rootPath string, info ServerInfo) error {
	pid := os.Getpid()

	c, err := newClient(endPoints)
	if err != nil {
		return err
	}

	keysAPI := client.NewKeysAPI(c)

	s := Service{
		ProcessId: pid,
		RootPath:  rootPath,
		Info:      info,
		KeysAPI:   keysAPI,
	}

	go s.HeartBeat()
	return nil
}

func (s *Service) HeartBeat() {
	api := s.KeysAPI
	serviceId := s.Info.String()

	for {
		key := s.RootPath + "/service_" + serviceId
		value, _ := json.Marshal(s.Info)

		_, err := api.Set(context.Background(), key, string(value), &client.SetOptions{
			TTL: time.Second * 60, //设置该 key TTL 为60秒,类似redis的expire time
		})

		if err != nil {
			xlog.Errorf("update [%s] error:%s\n", serviceId, err.Error())
		}
		time.Sleep(time.Second * 15)
	}
}
