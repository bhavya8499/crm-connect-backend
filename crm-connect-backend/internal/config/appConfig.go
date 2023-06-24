package config

import (
	"github.com/Netflix/go-env"
	"os"
	"strconv"
	"sync"
)

var (
	appConfig     Config
	onceAppConfig sync.Once
)

type Config struct {
	Env              string `env:"GO_ENV,default=local"`
	LogLevel         string `env:"LOG_LEVEL,default=debug"`
	Port             string `env:"PORT,default=4000"`
	ServiceName      string `env:"SERVICE_NAME,default=crm-connect-backend"`
	MiddlewareConfig *MiddlewareConfig
	ServerConfig     *ServerConfig
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// get environment variable as int
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return valueAsInt
}

// get environment variable as bool
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	valueAsBool, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return valueAsBool
}

// get environment variable as int64
func getEnvAsInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	valueAsInt64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return valueAsInt64
}

func NewAppConfig() *Config {

	onceAppConfig.Do(func() {
		_, _ = env.UnmarshalFromEnviron(&appConfig)
		appConfig.ServerConfig = NewServerConfig()
		appConfig.MiddlewareConfig = NewMiddlewareConfig()
	})
	return &appConfig
}
