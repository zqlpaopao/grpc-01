package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	v11 "grpc/test/src/proto"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func main(){
	RunServer(context.Background(),"9001")
}

func RunServer(ctx context.Context, port string) error {

	certFile := "/Users/zhangsan/Documents/GitHub/grpc-01/code/config/server.pem"
	keyFile := "/Users/zhangsan/Documents/GitHub/grpc-01/code/config/server.key"

	cTls := Tls()
	opts := []grpc.ServerOption{
		cTls,
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	mux := GetHTTPServeMux()

	server := grpc.NewServer(opts...)
	v11.RegisterSayHelloServiceServer(server, NewSayHelloResponseService())
	http.ListenAndServeTLS(":"+port,
		certFile,
		keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
			return
		}),
	)

	return nil

}
func GetHTTPServeMux() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testhttp: test-http-grpc"))
	})
	return mux
}

type Services struct {

}

func NewSayHelloResponseService()*Services{
	return &Services{}
}



func(s *Services) SayHello(ctx context.Context,req *v11.SayHelloRequest)(resp *v11.SayHelloResponse,err error){
	return &v11.SayHelloResponse{
		Response:             "http resp",
	}, err
}

/**************************************** 获取证书 **********************************/
func Tls()(grpc.ServerOption){
	c, err := credentials.NewServerTLSFromFile("/Users/zhangsan/Documents/GitHub/grpc-01/code/config/server.pem", "/Users/zhangsan/Documents/GitHub/grpc-01/code/config/server.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
		return nil
	}
	return grpc.Creds(c)
}





/**************************************** 拦截器 **********************************/
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
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


