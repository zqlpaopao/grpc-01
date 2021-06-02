package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc/test/src/proto"
	"io"
	"log"
)
const (
	PORT = "9002"
)
func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := pb.NewStreamServiceClient(conn)

	err = printSleep(client, &pb.PublicRequest{
		Req: &pb.Item{
			Value:                "valueSleep",
			Value2:               "value1Sleep",
		},
	})
	if err != nil {
		log.Fatalf("printSleep.err: %v", err)
	}
}

//双向流
func printSleep(client pb.StreamServiceClient, r *pb.PublicRequest) error {
	stream, err := client.Sleep(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: value1: %s, value2: %s", resp.Resp.Value, resp.Resp.Value2)
	}
	if err = stream.CloseSend();nil != err{
		log.Println(err)
	}
	return nil
}
