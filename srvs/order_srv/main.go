package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	"go-shop/srvs/order_srv/handler"
	"go-shop/srvs/order_srv/utils/otgrpc"
	"go-shop/srvs/order_srv/utils/register/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"go-shop/srvs/order_srv/global"
	"go-shop/srvs/order_srv/initialize"
	"go-shop/srvs/order_srv/proto"
	"go-shop/srvs/order_srv/utils"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "Port")

	//int
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvConn()
	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	//初始化jaeger
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.0.104:6831",
		},
		ServiceName: "mxshop",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))

	proto.RegisterOrderServer(server, &handler.OrderServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//註冊服務健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//啟動服務
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//服務註冊
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服務註冊失敗:", err.Error())
	}
	zap.S().Debugf("啟動服務器, 端口： %d", *Port)

	//監聽訂單超時topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.0.1:9876"}),
		consumer.WithGroupName("mxshop-order"),
	)

	if err := c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		fmt.Println("讀取消息失敗")
	}
	_ = c.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	_ = closer.Close()
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("註銷失敗:", err.Error())
	} else {
		zap.S().Info("註銷成功:")
	}
}
