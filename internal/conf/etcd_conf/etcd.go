package etcd_conf

import (
    "github.com/go-kratos/kratos/contrib/config/etcd/v2"
    "github.com/go-kratos/kratos/v2/config"
    clientv3 "go.etcd.io/etcd/client/v3"
    "google.golang.org/grpc"
    "time"
)

type etcdNode struct {
    path string
    clientNode *clientv3.Client
    s config.Source
}

func NewEtcSource(p string, serverList []string, tmout time.Duration) config.Source {
    client, err := clientv3.New(clientv3.Config{Endpoints: serverList,
        DialTimeout: tmout, DialOptions: []grpc.DialOption{grpc.WithBlock()}})
    if err != nil  || client == nil {
        return nil
    }
    ret, err := etcd.New(client, etcd.WithPath(p))
    if err != nil  {
        return nil
    }
    return &etcdNode{
        s: ret,
        path: p,
        clientNode: client,
    }
}

func (e *etcdNode) Load() ([]*config.KeyValue, error) {
    kvs, err :=  e.s.Load()
    if err != nil {
        return nil, err
    }
    /**
    for i, _ := range kvs {
        kvs[i].Key = ""
    }

     */
    return kvs, err
}

func (e *etcdNode) Watch() (config.Watcher, error) {
    return e.s.Watch()
}
