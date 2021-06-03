package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v11 "grpc/test/src/proto"
	"log"
	"net"
	"os"
	"runtime/debug"
)

func main(){
	RunServer(context.Background(),"9001")
}

func RunServer(ctx context.Context, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if nil != err {
		return err
	}

	/*
	var (
		opts []grpc.ServerOption
	)

	注册日志，各种时间因素
	opts = append(opts, gRpcServices.RegisterLogInject(logLayout, constant.GRpcLoginInsKey))
	opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 15 * time.Second, //client空闲超过该时间，发送一个GOAWAY
		//MaxConnectionAge:      time.Duration(math.MaxInt64), //client最大存活时间
		MaxConnectionAge:      5 * time.Second, //client最大存活时间
		MaxConnectionAgeGrace: 5 * time.Second, //强制关闭连接前缓冲时间，用以完成pending的请求
		Time:                  5 * time.Second, //client空闲该时间侯，发送一个ping
		Timeout:               3 * time.Second, //如果ping该时间内未收到pong，认为连接已断开
	}))
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             3 * time.Second, //client两次ping最小间隔，小于该时间中止连接
		PermitWithoutStream: true,            //即使没有活动的stream，也允许keepalive的ping
	}))

	 */
	//sv = grpc.NewServer(opts...)

	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}


	server := grpc.NewServer(opts...)
	v11.RegisterSayHelloServiceServer(server, NewSayHelloResponseService())
	c := make(chan os.Signal, 1)
	go func() {
		for range c {
			log.Println("shutting down GRPC server...")
			server.GracefulStop()//平滑关闭服务
			<-ctx.Done()
		}
	}()
	log.Println("start gRPC server...,port " + port)
	return server.Serve(listen)

}


type Services struct {

}

func NewSayHelloResponseService()*Services{
	return &Services{}
}



func(s *Services) SayHello(ctx context.Context,req *v11.SayHelloRequest)(resp *v11.SayHelloResponse,err error){
	return &v11.SayHelloResponse{
		Response:             "resp",
	}, err
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}


