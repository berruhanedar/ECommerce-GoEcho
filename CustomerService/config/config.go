package config

import "fmt"

type DbConfig struct {
	DBName  string
	ColName string
}

type CustomerConfig struct {
	Port     string
	DbConfig DbConfig
}

var RoleStatus = struct {
	System struct {
		NonPremium string
		Premium    string
	}
	Membership struct {
		Admin   string
		Manager string
		User    string
	}
}{
	System: struct {
		NonPremium string
		Premium    string
	}{
		NonPremium: "non-premium",
		Premium:    "premium",
	},
	Membership: struct {
		Admin   string
		Manager string
		User    string
	}{
		Admin:   "admin",
		Manager: "manager",
		User:    "user",
	},
}

var commonConfig = CustomerConfig{
	Port: ":8001",
	DbConfig: DbConfig{
		DBName:  "tesodev",
		ColName: "customer",
	},
}

var cfgs = map[string]CustomerConfig{
	"prod": commonConfig,
	"qa":   commonConfig,
	"dev":  commonConfig,
}

func GetCustomerConfig(env string) *CustomerConfig {
	config, exists := cfgs[env]
	if !exists {
		panic(fmt.Sprintf("config for env %s does not exist", env))
	}
	return &config
}
