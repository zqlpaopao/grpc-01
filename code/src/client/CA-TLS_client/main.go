package main
import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc/test/src/proto"
)
const PORT = "9002"
func main() {
	cert, err := tls.LoadX509KeyPair("/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/client/client.pem", "/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/client/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "test-grpc",
		RootCAs:      certPool,
	})


	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
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
