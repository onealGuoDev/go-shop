package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go-shop/api/good-api/global"
	"go-shop/api/good-api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	goodConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		zap.S().Fatal("連結使用者服務失敗")
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodConn)
}
