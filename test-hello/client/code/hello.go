package main

import (
	pb "github.com/freewebsys/grpc-go-demo/src/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	/*
	* 1）创建一个客户端连接 conn
	* 2）通过一个 conn 创建一个客户端
	* 3）发起 rpc 调用
	 */

	/*
		1、连创建一个客户端连接 conn  底层本质上调用newHTTP2Client，与server建立http2连接
		创建 conn连接 通过 Dial 方法创建 conn，Dial 调用了 DialContext 方法
		跟进 DialContext，发现 DialContext
		具体就是先实例化了一个 ClientConn 的结构体
		然后主要为 ClientConn 的 dopts 的各个属性进行初始化赋值
	*/

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//关闭连接
	defer conn.Close()
	/*
		创建客户端存根对象(此次请求的客户端初始化信息)，调用的这个方法在生成的pb.go文件中
		创建一个 greeterClient 的客户端
		主要参数1
			csMgr        *connectivityStateManager
			初始化 连接的状态管理器，每个连接具有 “IDLE”、“CONNECTING”、“READY”、“TRANSIENT_FAILURE”、“SHUTDOW N”、“Invalid-State” 这几种状态。
		主要参数2
			blockingpicker    *pickerWrapper
			pickerWrapper 是对 balancer.Picker 的一层封装，balancer.Picker 其实是一个负载均衡器，
			它里面只有一个 Pick 方法，它返回一个 SubConn 连接。

			SubConn
			分布式环境下，可能会存在多个 client 和 多个 server，client 发起一个 rpc 调用之前，
			需要通过 balancer 去找到一个 server 的 address，balancer 的 Picker 类返回一个 SubConn，
			SubConn 里面包含了多个 server 的 address，假如返回的 SubConn 是 “READY” 状态，grpc 会发送 RPC 请求，
			否则则会阻塞，等待 UpdateBalancerState 这个方法更新连接的状态并且通过 picker 获取一个新的 SubConn 连接。

	*/
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	/*
		发起调用
		SayHello 方法是通过调用 Invoke 的方法去发起 rpc 调用
		Invoke 方法调用了 invoke， 在 invoke
		这个方法里面 sendMsg recvMsg 两个接口 这两个接口在 clientStream 中被实现了。
		sendMsg 发送消息====================================
			先准备数据，然后再调用 csAttempt 这个结构体中的 sendMsg 方法
			最终是通过 a.t.Write 发出的数据写操作 a.t 是一个 ClientTransport 类型，
			所以最终是通过 ClientTransport 这个结构体的 Write 方法发送数据
		recvMsg  接收消息================================
			调用了 recv 方法  最终落在了
			p.r.Read 方法，p.r 是一个 io.Reader 类型
			最终都是要落到 IO
	*/

	ctx, connel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer connel()
	//time.Sleep(10 * time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(1*time.Second)))
	//defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("####### get server Greeting response: %s", r.Message)
}
