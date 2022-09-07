package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/satori/go.uuid"
	"go-shop/srvs/inventory_srv/handler"
	"go-shop/srvs/inventory_srv/utils/register/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-shop/srvs/inventory_srv/global"
	"go-shop/srvs/inventory_srv/initialize"
	"go-shop/srvs/inventory_srv/proto"
	"go-shop/srvs/inventory_srv/utils"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 8003, "Port")

	//init
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
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
		zap.S().Panic("服務註冊失败:", err.Error())
	}
	zap.S().Debugf("啟動服務器, 端口： %d", *Port)

	//監視庫存歸還的topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.0.1:9876"}),
		consumer.WithGroupName("mxshop-inventory"),
	)

	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("讀取消息失败")
	}
	_ = c.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("註銷失败:", err.Error())
	} else {
		zap.S().Info("註銷成功:")
	}
}
