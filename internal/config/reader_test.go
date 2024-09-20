package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/WsDev69/producer-consumer/internal/config"
)

func TestReadConfig(t *testing.T) {
	// Set environment variables

	t.Setenv("PREFIX_LOGLEVEL", "info")
	t.Setenv("PREFIX_OUTPUT", "plain")
	t.Setenv("PREFIX_MESSAGERATE", "200")
	t.Setenv("PREFIX_POSTGRES_HOST", "localhost")
	t.Setenv("PREFIX_POSTGRES_PORT", "5432")
	t.Setenv("PREFIX_POSTGRES_PASSWORD", "password123")
	t.Setenv("PREFIX_POSTGRES_USER", "testuser")
	t.Setenv("PREFIX_POSTGRES_DB", "testdb")
	t.Setenv("PREFIX_POSTGRES_MAXIDLECONN", "15")
	t.Setenv("PREFIX_GRPC_HOST", "localhost")
	t.Setenv("PREFIX_GRPC_PORT", "5000")
	t.Setenv("PREFIX_PROMETHEUS_PORT", "9090")

	// Read configuration using the Read function
	cfg, err := config.Read("PREFIX")
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Assertions to check if the values are set correctly
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "plain", cfg.OutPut)
	assert.Equal(t, 200, cfg.MessageRate)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
	assert.Equal(t, "password123", cfg.Postgres.Password)
	assert.Equal(t, "testuser", cfg.Postgres.User)
	assert.Equal(t, "testdb", cfg.Postgres.DB)
	assert.Equal(t, 15, cfg.Postgres.MaxIdleConn)
	assert.Equal(t, "localhost", cfg.GRPC.Host)
	assert.Equal(t, 5000, cfg.GRPC.Port)
	assert.Equal(t, 9090, cfg.Prometheus.Port)
}

func TestReadProducerConfig(t *testing.T) {
	// Set environment variables for Producer
	t.Setenv("PRODUCER_LOGLEVEL", "debug")
	t.Setenv("PRODUCER_NUMWORKER", "20")
	t.Setenv("PRODUCER_MAXBACKLOG", "1000")

	// Read producer config using the ReadProducer function
	producerCfg, err := config.ReadProducer("PRODUCER")
	require.NoError(t, err)
	assert.NotNil(t, producerCfg)

	// Assertions
	assert.Equal(t, "debug", producerCfg.LogLevel)
	assert.Equal(t, 20, producerCfg.NumWorker)
	assert.Equal(t, int64(1000), producerCfg.MaxBacklog)
}

func TestReadConsumerConfig(t *testing.T) {
	// Set environment variables for Consumer
	t.Setenv("CONSUMER_LOGLEVEL", "warn")
	t.Setenv("CONSUMER_MESSAGERATE", "150")

	// Read consumer config using the ReadConsumer function
	consumerCfg, err := config.ReadConsumer("CONSUMER")
	require.NoError(t, err)
	assert.NotNil(t, consumerCfg)

	// Assertions
	assert.Equal(t, "warn", consumerCfg.LogLevel)
	assert.Equal(t, 150, consumerCfg.MessageRate)
}
