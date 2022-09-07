package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go-shop/api/user_api/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go-shop/api/user_api/global"
	"go-shop/api/user_api/initialize"
	myvalidator "go-shop/api/user_api/validator"
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
	if !debug {
		//port, err := utils.GetFreePort()
		//if err == nil {
		//	global.ServerConfig.Port = port
		//}
	}

	//
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "格式錯誤的手機!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服務註冊失敗:", err.Error())
	}

	zap.S().Debugf("啟動服務器, Port： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("啟動失敗:", err.Error())
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("註銷失敗:", err.Error())
	} else {
		zap.S().Info("註銷成功:")
	}
}
