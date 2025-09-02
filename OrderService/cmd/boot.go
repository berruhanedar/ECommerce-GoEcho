package cmd

import (
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/client"
	"tesodev-korpes/pkg/middleware"
	"tesodev-korpes/shared/config"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func BootOrderService(clientMongo *mongo.Client, e *echo.Echo) {
	cfg := config3.GetOrderConfig("dev")

	orderCol, err := pkg.GetMongoCollection(clientMongo, cfg.DbConfig.DBName, cfg.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(orderCol)
	cc := client.New("http://localhost:8001", 5*time.Second)
	service := internal.NewService(repo, cc)
	handler := internal.NewHandler(e, service, clientMongo)


	e.GET("/price/:id",
		func(c echo.Context) error {
			e.Router().Find(c.Request().Method, c.Path(), c)
			return c.Handler()(c)
		},
		middleware.Authentication(clientMongo, nil),
		middleware.AuthorizationMiddleware(config.Cfg.AllowedRoles),
		middleware.RoleRouting(config.Cfg),
	)
	e.GET("/internal/price/premium/:id", handler.GetPremiumOrderPrice)
	e.GET("/internal/price/non-premium/:id", handler.GetNonPremiumOrderPrice)

	e.Logger.Fatal(e.Start(cfg.Port))
}
