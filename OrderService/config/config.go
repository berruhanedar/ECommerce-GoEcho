package config

type OrderConfig struct {
	Port     string
	DbConfig struct {
		DBName  string
		ColName string
	}
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

var cfgs = map[string]OrderConfig{
	"prod": {
		Port: ":8002",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "order",
		},
	},
	"qa": {
		Port: ":8002",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "order",
		},
	},
	"dev": {
		Port: ":8002",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "order",
		},
	},
}

func GetOrderConfig(env string) *OrderConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
