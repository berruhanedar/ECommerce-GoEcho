package main

import (
	"fmt"
	customercmd "tesodev-korpes/CustomerService/cmd"
	ordercmd "tesodev-korpes/OrderService/cmd"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/middleware"
	"tesodev-korpes/shared/config"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	dbConf := config.GetDBConfig("dev")

	client, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		panic(err)
	}
	fmt.Println("connecting db")

	customerEcho := echo.New()
	customerEcho.Use(middleware.CorrelationIdMiddleware())
	customerEcho.Use(middleware.LoggingMiddleware)
	customerEcho.Use(middleware.RecoveryMiddleware)
	customerEcho.Use(middleware.ErrorHandler())
	customerEcho.GET("/swagger/*", echoSwagger.WrapHandler)

	orderEcho := echo.New()
	orderEcho.Use(middleware.CorrelationIdMiddleware())
	orderEcho.Use(middleware.LoggingMiddleware)
	orderEcho.Use(middleware.RecoveryMiddleware)
	orderEcho.Use(middleware.ErrorHandler())
	orderEcho.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		customercmd.BootCustomerService(client, customerEcho)
	}()

	ordercmd.BootOrderService(client, orderEcho)
}
