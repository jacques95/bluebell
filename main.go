package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	// 1、 	加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("settings initialize failed, err:%v\n", err)
		return
	}

	// 2、 	初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger initialize failed, err:%v\n", err)
		return
	}
	defer func(l *zap.Logger) {
		if err := l.Sync(); err != nil {
			return
		}
	}(zap.L())

	// 3、 	初始化 mysql 连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("mysql initialize failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 4、 	初始化 redis 连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis initialize failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	//	5、初始化 ID 生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	//	6、初始化 gin.validator i18
	//if err := controller.InitTrans("zh"); err != nil {
	//	fmt.Printf("init validator trans failed, err:%v\n", err)
	//	return
	//}

	//	7、注册路由
	r := routes.Setup()

	//	8、启动服务
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d",
			viper.GetInt("port"),
		),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown.", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
