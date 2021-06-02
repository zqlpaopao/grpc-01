package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "grpc/test/src/proto"
	"io"
	"log"
	"net"
)
type StreamService struct{}

const (
	PORT = "9002"
)
func main() {

	//设置客户端最大接收值
	//opt := grpc.MaxRecvMsgSize()
	//grpc.MaxSendMsgSize()

	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}
//客户端流事rpc
func (s *StreamService) Work(stream pb.StreamService_WorkServer) error {

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
	for {
		r ,err := stream.Recv()
		if err == io.EOF{
			return stream.SendAndClose(&pb.PublicResponse{
				Resp:                &pb.Item{
					Value:                "client-stream-server",
					Value2:               "client-stream-server-v2",
				} ,
			})
		}
		if err != nil{
			return err
		}
		log.Printf("stream.Recv value: %s,value2: %s", r.Req.Value, r.Req.Value2)
	}
}

//服务端流式
func (s *StreamService) Eat(r *pb.PublicRequest, stream pb.StreamService_EatServer) error {

	return nil
}

func (s *StreamService) Sleep(stream pb.StreamService_SleepServer) error {
	return nil
}

