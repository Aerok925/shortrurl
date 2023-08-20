package configs

import "flag"

type Config struct {
	Server Server
	Result Result
}

type Server struct {
	Address string
}

type Result struct {
	BaseAddress string
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
	return &Config{
		Server: Server{
			Address: *address,
		},
		Result: Result{
			BaseAddress: *baseAddress,
		},
	}
}
