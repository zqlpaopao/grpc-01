package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	v11 "grpc/test/src/proto"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"log"
	"net"
	"os"
)

const (
	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)


func main(){
	{
		// set up a span reporter
		reporter := zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
		defer reporter.Close()

		// create our local service endpoint
		endpoint, err := zipkin.NewEndpoint(SERVICE_NAME, ZIPKIN_RECORDER_HOST_PORT)
		if err != nil {
			log.Fatalf("unable to create local endpoint: %+v\n", err)
		}

		// initialize our tracer
		nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
		if err != nil {
			log.Fatalf("unable to create tracer: %+v\n", err)
		}

		// use zipkin-go-opentracing to wrap our tracer
		tracer := zipkinot.Wrap(nativeTracer)

		// optionally set as Global OpenTracing tracer instance
		opentracing.SetGlobalTracer(tracer)
	}

	RunServer(context.Background(),"9001")
}

func RunServer(ctx context.Context, port string) error {

	cs, err := credentials.NewServerTLSFromFile("/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/server/server.pem", "/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/server/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	server := grpc.NewServer(grpc.Creds(cs))
	listen, err := net.Listen("tcp", ":"+port)
	if nil != err {
		return err
	}


	v11.RegisterSayHelloServiceServer(server, NewSayHelloResponseService())
	c := make(chan os.Signal, 1)
	go func() {
		for range c {
			log.Println("shutting down GRPC server...")
			server.GracefulStop()//平滑关闭服务
			<-ctx.Done()
		}
	}()
	log.Println("start gRPC server...,port " + port)
	return server.Serve(listen)

}

type Auth struct {
	appKey    string
	appSecret string
}

type Services struct {
}

func NewSayHelloResponseService()*Services{
	return &Services{}
}



func(s *Services) SayHello(ctx context.Context,req *v11.SayHelloRequest)(resp *v11.SayHelloResponse,err error){
	return &v11.SayHelloResponse{
		Response:             "resp",
	}, err
}


