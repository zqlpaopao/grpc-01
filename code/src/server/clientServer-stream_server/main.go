package main

import (
	"google.golang.org/grpc"
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
	return nil
}

//服务端流式
func (s *StreamService) Eat(r *pb.PublicRequest, stream pb.StreamService_EatServer) error {

	return nil
}

//双向流
func (s *StreamService) Sleep(stream pb.StreamService_SleepServer) error {
	n := 0
	for {
		err := stream.Send(&pb.PublicResponse{
			Resp: &pb.Item{
				Value:  "gPRC Stream Client: Sleep",
				Value2: "双向stream-value2",
			},
		})
		if err != nil {
			return err
		}
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("stream.Recv req.value: %s, pt.value2: %s", r.Req.Value, r.Req.Value2)
	}

}

