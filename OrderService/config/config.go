package config

type DbConfig struct {
	DBName  string
	ColName string
}

type OrderConfig struct {
	Port     string
	DbConfig DbConfig
}

var OrderStatus = struct {
	Ordered   string
	Shipped   string
	Delivered string
	Canceled  string
}{
	Ordered:   "ORDERED",
	Shipped:   "SHIPPED",
	Delivered: "DELIVERED",
	Canceled:  "CANCELED",
}

var commonConfig = OrderConfig{
	Port: ":8002",
	DbConfig: DbConfig{
		DBName:  "tesodev",
		ColName: "order",
	},
}

var cfgs = map[string]OrderConfig{
	"prod": commonConfig,
	"qa":   commonConfig,
	"dev":  commonConfig,
}

func GetOrderConfig(env string) *OrderConfig {
	cfg, exists := cfgs[env]
	if !exists {
		panic("config does not exist")
	}
	return &cfg
}
