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

	err = printEats(client, &pb.PublicRequest{
		Req:                  &pb.Item{
			Value:                "value",
			Value2:               "value1",

		},
	})
	if err != nil {
		log.Fatalf("printEats.err: %v", err)
	}

	err = printWork(client, &pb.PublicRequest{
		Req: &pb.Item{
			Value:                "valueWork",
			Value2:               "value1Work",

		},
	})
	if err != nil {
		log.Fatalf("printWork.err: %v", err)
	}
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

//服务端流式
func printEats(client pb.StreamServiceClient, r *pb.PublicRequest) error {
	var c context.Context

	c = context.WithValue(context.TODO(),"a","b")

	stream,err := client.Eat(c,r)
	if err != nil{
		return err
	}
	//接收server的header信息
	fmt.Println(stream.Header())//map[cc:[dd] content-type:[application/grpc]]

	for {
		resp ,err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}
		log.Printf("resp: value1 %s, value1 %s",resp.Resp.Value,resp.Resp.Value2)
	}
	//在一元rpc中header和trailer是一起到达的，在流式中是在接受消息后到达的
	fmt.Println(stream.Trailer())//map[cc1:[dd1]]

	return nil
}
func printWork(client pb.StreamServiceClient, r *pb.PublicRequest) error {
	return nil
}
func printSleep(client pb.StreamServiceClient, r *pb.PublicRequest) error {
	return nil
}
