package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	v11 "grpc/test/src/proto"
	"log"
	"os"
)
const PORT = "9001"

func main(){

	c := Tls()
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()


	client := v11.NewSayHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &v11.SayHelloRequest{
		Request: "gRPC-http",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}



func Tls()(credentials.TransportCredentials){

	c, err := credentials.NewClientTLSFromFile("/Users/zhangsan/Documents/GitHub/grpc-01/code/config/server.pem", "test-http-grpc")
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
		return nil
	}

	return c
}

