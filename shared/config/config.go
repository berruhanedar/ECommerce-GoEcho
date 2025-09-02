package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	MongoDuration  time.Duration
	MongoClientURI string
}

type ServiceUrls struct {
	CustomerServiceURL string
}

type AuthConfig struct {
	JWTSecret string
}

type Config struct {
	RoleMapping  map[string]string
	AllowedRoles []string
}

var Cfg = Config{
	AllowedRoles: []string{"premium", "non-premium"},
	RoleMapping: map[string]string{
		"premium":     "/internal/price/premium/:id",
		"non-premium": "/internal/price/non-premium/:id",
	},
}

var cfgs = map[string]DbConfig{
	"prod": {
		MongoDuration: time.Second * 100,
	},
	"qa": {
		MongoDuration: time.Second * 100,
	},
	"dev": {
		MongoDuration: time.Second * 100,
	},
}

var (
	serviceUrls ServiceUrls
	authConfig  AuthConfig
)

func init() {

	if err := godotenv.Load("./media/.env"); err != nil {
		panic("Environment variable did not load")
	}
	fmt.Println("Environment variables loaded")

	serviceUrls = ServiceUrls{
		CustomerServiceURL: os.Getenv("CUSTOMER_SERVICE_URL"),
	}
	authConfig = AuthConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if authConfig.JWTSecret == "" {
		panic("JWT_SECRET environment variable not set")
	}
}

func GetDBConfig(env string) *DbConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}

	config.MongoClientURI = os.Getenv("MONGO_URI")
	return &config
}
func GetServiceURLs() ServiceUrls {
	return serviceUrls
}

func GetAuthConfig() AuthConfig {
	return authConfig
}
