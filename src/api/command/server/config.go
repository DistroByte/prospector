package server

import (
	"net"
	"strconv"
)

type Config struct {
	BindAddr string
	Port     int
}

func DefaultConfig() *Config {
	cfg := &Config{
		BindAddr: "0.0.0.0",
		Port:     3434,
	}

	return cfg
}

func (c *Config) Merge(b *Config) *Config {
	result := *c

	if b.BindAddr != "" {
		result.BindAddr = b.BindAddr
	}

	if b.Port != 0 {
		result.Port = b.Port
	}

	return &result
}

func (c *Config) Listener(proto, addr string, port int) (net.Listener, error) {
	if addr == "" {
		addr = c.BindAddr
	}

	return net.Listen(proto, net.JoinHostPort(addr, strconv.Itoa(port)))
}
