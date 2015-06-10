package etcd

import (
	"log"

	client "github.com/coreos/go-etcd/etcd"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
}

// Client() returns a new client for accessing etcd.
//
func (cfg *Config) Client() (*client.Client, error) {

	log.Printf("[INFO] Consul etcd configured with endpoints: %s", cfg.Endpoint)
	c := client.NewClient([]string{cfg.Endpoint})

	// should we "ping" here to ensure we can communicate?

	return c, nil
}
