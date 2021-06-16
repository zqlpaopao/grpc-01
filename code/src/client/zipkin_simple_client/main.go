package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	v11 "grpc/test/src/proto"
	"log"
)
const (
	PORT = "9001"
	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)


func main() {
	c, err := credentials.NewClientTLSFromFile("/Users/zhangsan/Documents/GitHub/grpc-01/code/conf/server/server.pem", "test-grpc")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	//{
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
	//}

	conn, err := grpc.Dial(":"+PORT,
		grpc.WithTransportCredentials(c),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
		))
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

	resp ,err = client.SayHello(context.Background(),&v11.SayHelloRequest{
		Request:              "hello",

	})
	fmt.Println(resp)
	fmt.Println(err)
	//os.Exit(3)
}

