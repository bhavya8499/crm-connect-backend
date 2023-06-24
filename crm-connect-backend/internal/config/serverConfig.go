package config

import (
	"github.com/Netflix/go-env"
	"sync"
)

type ServerConfig struct {
	HTTPPort                  string `env:"HTTP_PORT,default=8080"`
	HTTPServerShutDownTimeout int    `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT,default=5"`
	GRPCPort                  string `env:"GRPC_PORT,default=80"`
}

var (
	serverConfig     ServerConfig
	onceServerConfig sync.Once
)

// NewServerConfig is singleton and it makes the configs exportable
func NewServerConfig() *ServerConfig {
	onceServerConfig.Do(func() {
		_, _ = env.UnmarshalFromEnviron(&serverConfig)
	})
	return &serverConfig
}
