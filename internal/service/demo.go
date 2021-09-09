package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	pb "helloworld/api/helloworld/v1"
)

type DemoService struct {
	pb.UnimplementedDemoServer
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

func (s *DemoService) CreateDemo(ctx context.Context, req *pb.CreateDemoRequest) (*pb.CreateDemoReply, error) {
	return &pb.CreateDemoReply{}, nil
}
func (s *DemoService) UpdateDemo(ctx context.Context, req *pb.UpdateDemoRequest) (*pb.UpdateDemoReply, error) {
	return &pb.UpdateDemoReply{}, nil
}
func (s *DemoService) DeleteDemo(ctx context.Context, req *pb.DeleteDemoRequest) (*pb.DeleteDemoReply, error) {
	return &pb.DeleteDemoReply{}, nil
}
func (s *DemoService) GetDemo(ctx context.Context, req *pb.GetDemoRequest) (*pb.GetDemoReply, error) {
	return &pb.GetDemoReply{}, nil
}
func (s *DemoService) ListDemo(ctx context.Context, req *pb.ListDemoRequest) (*pb.ListDemoReply, error) {
	return &pb.ListDemoReply{}, nil
}
func (s *DemoService) DemoReq(ctx context.Context, req *pb.CreateDemoRequest) (*pb.CreateDemoReply, error) {
	/**
	//这是对接etcd 服务注册和发现
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"0.0.0.0:2379"},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(cli)
	 */
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(cli)

	conn, err := transgrpc.DialInsecure(
		context.Background(),
		//transgrpc.WithEndpoint("127.0.0.1:9000"), //使用创建客户端连接
		transgrpc.WithEndpoint("discovery:///hello_world"), // 使用服务发现
		transgrpc.WithDiscovery(r),  // 使用服务发现
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		), )
	/*
	//这是创建 调用http client conn
	conn, err := http.NewClient(
			context.Background(),
			http.WithMiddleware(
				recovery.Recovery(),
			),
			http.WithEndpoint("discovery:///hello_world"),
			http.WithDiscovery(r),
		)
	 */

	if err != nil {
		return &pb.CreateDemoReply{
			Data: fmt.Errorf("new client err: %v", err).Error(),
		}, nil
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	/*
		client := NewGreeterHTTPClient(conn) // 这是创建http client
	 */
	rsp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name:"rpc test name", Data:"rpc test data"})
	if err != nil {
		return &pb.CreateDemoReply{
			Data: fmt.Errorf("call sayHello err: %v", err).Error(),
		}, nil
	}
	return  &pb.CreateDemoReply{
		Data: fmt.Sprintf("say rpc hello value: %s", rsp.Message),
	}, nil
}
