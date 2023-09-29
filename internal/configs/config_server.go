package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Server Server
	Result Result
}

type Server struct {
	Address string `env:"SERVER_ADDRESS"`
}

type Result struct {
	BaseAddress string `env:"BASE_URL"`
}

var (
	address     *string
	baseAddress *string
)

func init() {
	address = flag.String("a", "localhost:8080", "server address")
	baseAddress = flag.String("b", "http://localhost:8080", "base address")
}

func New() *Config {
	flag.Parse()
	cfg := &Config{
		Server: Server{},
		Result: Result{},
	}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	if cfg.Server.Address == "" {
		cfg.Server.Address = *address
	}
	if cfg.Result.BaseAddress == "" {
		cfg.Result.BaseAddress = *baseAddress
	}
	log.Println(cfg)
	return cfg
}
