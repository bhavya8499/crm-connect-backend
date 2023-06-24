package config

import (
	"github.com/Netflix/go-env"
	"strings"
	"sync"
)

type MiddlewareConfig struct {
	Tokens        string `env:"TOKENS,default=39852a79-0c49-40f6-9dec-0247048d631d"`
	EnabledTokens []string
}

var (
	middlewareConfig     MiddlewareConfig
	onceMiddlewareConfig sync.Once
)

func NewMiddlewareConfig() *MiddlewareConfig {
	onceMiddlewareConfig.Do(func() {
		_, _ = env.UnmarshalFromEnviron(&middlewareConfig)
		middlewareConfig.EnabledTokens = strings.Split(middlewareConfig.Tokens, ",")
	})
	return &middlewareConfig
}
