package config

type CustomerConfig struct {
	Port     string
	DbConfig struct {
		DBName  string
		ColName string
	}
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

var cfgs = map[string]CustomerConfig{
	"prod": {
		Port: ":8001",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "customer",
		},
	},
	"qa": {
		Port: ":8001",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "customer",
		},
	},
	"dev": {
		Port: ":8001",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "customer",
		},
	},
}

func GetCustomerConfig(env string) *CustomerConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
