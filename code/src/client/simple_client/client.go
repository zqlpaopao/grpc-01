/**
 * @Author: zhangSan
 * @Description:
 * @File:  client
 * @Version: 1.0.0
 * @Date: 2021/5/25 下午4:40
 */

package main

import (
	"context"
	"google.golang.org/grpc"
	v11 "grpc/test/src/proto"
	"log"
)
const PORT = "9001"

func main(){
		conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("grpc.Dial err: %v", err)
		}
		defer conn.Close()
		client := v11.NewSayHelloServiceClient(conn)
		resp, err := client.SayHello(context.Background(), &v11.SayHelloRequest{
			Request: "gRPC",
		})
		if err != nil {
			log.Fatalf("client.Search err: %v", err)
		}
		log.Printf("resp: %s", resp.GetResponse())
}


