package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go-shop/api/good-api/global"
	"go-shop/api/good-api/initialize"
	"go-shop/api/good-api/utils/register/consul"
)

func main() {
	//1. init logger
	initialize.InitLogger()

	//2. init 配置
	initialize.InitConfig()

	//3. init routers
	Router := initialize.Routers()

	//4. init 翻譯
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//5. 初始化srv的連結
	initialize.InitSrvConn()

	viper.AutomaticEnv()

	debug := viper.GetBool("GO_SHOP_DEBUG")
	//debug = false
	if !debug {
		//port, err := utils.GetFreePort()
		//if err == nil {
		//	global.ServerConfig.Port = port
		//}
	}

	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服務註冊失敗:", err.Error())
	}
	zap.S().Debugf("啟動服務器, Port： %d", global.ServerConfig.Port)
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("啟動失敗:", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("註銷失敗:", err.Error())
	} else {
		zap.S().Info("註銷成功:")
	}
}
