package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"helloworld/internal/conf/etcd_conf"
	"os"
	"time"

	// etcdclient "go.etcd.io/etcd/client/v3"
	// "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	c_api "github.com/hashicorp/consul/api"
	"helloworld/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	/**
	//这是etc服务发现
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"0.0.0.0:2379"},
	})

	if err != nil {
		return nil
	}
	r := etcd.New(client)
	 */

	consulClient, err := c_api.NewClient(c_api.DefaultConfig())
	if err != nil {
		return nil
	}
	r := consul.New(consulClient)

	Name = "hello_world"
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", log.TraceID(),
		"span_id", log.SpanID(),
	)


	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	//使用etcd作为配置中心:
	//defer client.Close()
	testKey := "/test_demo_etcd_config"
	e := config.New(
		config.WithSource(
			etcd_conf.NewEtcSource(testKey, []string{"127.0.0.1:2379"}, time.Second),
		),
	)
	if err := e.Load(); err != nil {
		panic(err)
	}

	var busCnf conf.Bootstrap
	var xyz struct {
		XY string `json:"xy"`
	}
	if err := e.Scan(&xyz); err != nil {
		panic(err)
	}
	fmt.Println(xyz)
	fmt.Println(busCnf)

	app, cleanup, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
