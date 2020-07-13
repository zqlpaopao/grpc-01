package main

import (
	"fmt"
	pb "github.com/freewebsys/grpc-go-demo/src/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const (
	port = ":50051"
)

//服务类
type server struct{}

// 服务方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(5 * time.Second)
	fmt.Println("######### get client request name :" + in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	/*
	*  初始化serve端参数默认值，头信息及
	*  对serverOptions（相当于http 的header信息） 进行初始化赋值为默认值
	*  初始化默认值----接收、发送消息，连接超时时间、读写缓冲区大小
	*  拦截器的加载一元和流式可以在服务端接收到请求时优先对请求中的数据做一些处理后再转交给指定的服务处理并响应，功能类似middleware
	*  被调用服务信息的追踪
	*  被调用服务的ListenSocket信息和方法的编号，原子性递增的，记录在服务注册中server(map)中
	 */
	s := grpc.NewServer()
	/*
	* pb(proto)生成的RegisterGreeterServer 调用初始化的server的 RegisterService
	* 通过反射获取方法名称和处理地址，并添加到server初始化的一元rpc 和流式rpc的服务列表中（map）
	* 检测注册的方法是否在服务列表中
	 */
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	//将注册的方法绑定到grpcserver中
	reflection.Register(s)
	/*
	* 启动服务端 accpet 接收请求
	* 如果accept 失败进行后续服务关闭 失败等待最长是1秒，监听失败进行错误返回，然后处理下一个
	*
	 */
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
