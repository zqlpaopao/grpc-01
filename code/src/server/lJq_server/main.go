package main
import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "grpc/test/src/proto"
	"runtime/debug"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)
type StreamService struct{}

const PORT = "9002"
func main() {
	cert, err := tls.LoadX509KeyPair("/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/server/server.pem", "/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/server/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	//拦截器
	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithStreamServerChain(
			//RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}

//客户端流rpc
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


//type StreamServerInterceptor func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error

func LoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Println("srv",srv)
	fmt.Printf("%#v\n",ss)
	fmt.Printf("%#v\n",info)

	err := handler(srv,ss)
	fmt.Println(err)

	return err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}