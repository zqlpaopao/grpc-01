/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/5/26 下午5:48
 */

package main

import (
	"context"
	"fmt"
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

	err = printWork(client, &pb.PublicRequest{
		Req: &pb.Item{
			Value:                "valueWork",
			Value2:               "value1Work",

		},
	})
	if err != nil {
		log.Fatalf("printWork.err: %v", err)
	}
}

func printWork(client pb.StreamServiceClient, r *pb.PublicRequest) error {
	stream,err := client.Work(context.Background())
	if err != nil{
		return err
	}

	for i := 0 ;i < 6;i++{
		fmt.Println(r)
		err := stream.Send(r)
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}
	}

	//注意这个header是设置不了的
	//fmt.Println(stream.Header())

	resp ,err := stream.CloseAndRecv()
	if err != nil{
		return err
	}

	log.Printf("resp: value1 %s, value1 %s",resp.Resp.Value,resp.Resp.Value2)

	//在一元rpc中header和trailer是一起到达的，在流式中是在接受消息后到达的
	fmt.Println(stream.Trailer())//map[cc1:[dd1]]
	return nil
}
