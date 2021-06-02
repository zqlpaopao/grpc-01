package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "grpc/test/src/proto"
	"log"
	"net"
	"time"
)
type StreamService struct{}

const (
	PORT = "9002"
)
func main() {

	//设置客户端最大接收值
	//opt := grpc.MaxRecvMsgSize()

	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", "127.0.0.1:"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}
//服务端流式
func (s *StreamService) Eat(r *pb.PublicRequest, stream pb.StreamService_EatServer) error {
	//设置header信息 sendHeader不可同时用，否则SendHeader会覆盖前一个
	if err := stream.SetHeader(metadata.MD{"cc2":[]string{"dd2"}});nil != err{
		return err
	}
	//设置header信息
	//if err := stream.SendHeader(metadata.MD{"cc":[]string{"dd"}});err != nil{
	//	return err
	//}



	//设置metadata，注意一元和流式的区别
	stream.SetTrailer(metadata.MD{"cc1":[]string{"dd1"}})

	a := stream.Context().Value("a")
	fmt.Println(a)

	for i := 0; i < 10;i++{
		time.Sleep(1 *time.Second)
		err := stream.Send(&pb.PublicResponse{
			Resp:                 &pb.Item{
				Value:                "eat",
				Value2:               "服务端流式",

			},
		})

		if err != nil{
			return err
		}



	}

	return nil
}

func (s *StreamService) Work(stream pb.StreamService_WorkServer) error {
	return nil
}
func (s *StreamService) Sleep(stream pb.StreamService_SleepServer) error {
	return nil
}
