package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"go-shop/srvs/user_srv/global"
	"go-shop/srvs/user_srv/handler"
	"go-shop/srvs/user_srv/initialize"
	"go-shop/srvs/user_srv/proto"
	"go-shop/srvs/user_srv/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 8001, "Port")

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
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//註冊服務健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服務註冊
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	fmt.Printf("%s", global.ServerConfig.ConsulInfo)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.0.100:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成註冊對象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"go_shop", "bobby", "user", "srv"}
	registration.Address = "192.168.0.100"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("註銷失敗")
	}
	zap.S().Info("註銷成功")
}
