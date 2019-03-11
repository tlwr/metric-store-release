package main

import (
	envstruct "code.cloudfoundry.org/go-envstruct"
	"github.com/cloudfoundry/metric-store/pkg/tls"
)

// Config is the configuration for a MetricStore.
type Config struct {
	LogProviderAddr string `env:"LOGS_PROVIDER_ADDR, required, report"`
	LogsProviderTLS LogsProviderTLS

	MetricStoreAddr       string `env:"METRIC_STORE_ADDR, required, report"`
	HealthAddr            string `env:"HEALTH_ADDR, report"`
	ShardId               string `env:"SHARD_ID, required, report"`
	TimerRollupBufferSize uint   `env:"TIMER_ROLLUP_BUFFER_SIZE, report"`
	NodeIndex             int    `env:"NODE_INDEX, report"`

	MetricStoreTLS tls.TLS
}

// LogsProviderTLS is the LogsProviderTLS configuration for a MetricStore.
type LogsProviderTLS struct {
	LogProviderCA   string `env:"LOGS_PROVIDER_CA_FILE_PATH, required, report"`
	LogProviderCert string `env:"LOGS_PROVIDER_CERT_FILE_PATH, required, report"`
	LogProviderKey  string `env:"LOGS_PROVIDER_KEY_FILE_PATH, required, report"`
}

// LoadConfig creates Config object from environment variables
func LoadConfig() (*Config, error) {
	c := Config{
		MetricStoreAddr:       ":8080",
		HealthAddr:            "localhost:6061",
		ShardId:               "metric-store",
		TimerRollupBufferSize: 16384,
	}

	if err := envstruct.Load(&c); err != nil {
		return nil, err
	}

	return &c, nil
}