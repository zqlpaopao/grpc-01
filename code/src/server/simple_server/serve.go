/**
 * @Author: zhangSan
 * @Description:
 * @File:  serve
 * @Version: 1.0.0
 * @Date: 2021/5/25 上午11:13
 */

package main

import (
	"context"
	"google.golang.org/grpc"
	v11 "grpc/test/src/proto"
	"log"
	"net"
	"os"
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


	server := grpc.NewServer()
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


