package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"producer-consumer/internal/config"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		name   string
		hasErr bool
		exp    *config.Config
		init   func(t *testing.T)
	}{
		{
			name:   "success",
			hasErr: false,
			exp: &config.Config{
				LogLevel: "info",
				Postgres: config.Postgres{
					Host:            "localhost",
					Port:            1234,
					Password:        "12345",
					User:            "postgres",
					DB:              "db",
					ConnMaxLifeTime: time.Minute * 5,
					MaxIdleConn:     10,
					MaxOpenConn:     10,
				},
				GRPC: config.GRPC{
					Host:     "10.10.8.70",
					Port:     8080,
					TLS:      false,
					CertFile: "",
					KeyFile:  "",
				},
			},
			init: func(t *testing.T) {
				t.Setenv("TEST_LOGLEVEL", "info")
				t.Setenv("TEST_POSTGRES_HOST", "localhost")
				t.Setenv("TEST_POSTGRES_PORT", "1234")
				t.Setenv("TEST_GRPC_HOST", "10.10.8.70")
			},
		},
		{
			name:   "err",
			hasErr: true,
			init: func(t *testing.T) {
				t.Setenv("TEST_POSTGRES_PORT", "port value")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.init(t)
			res, err := config.Read("test")
			if tc.hasErr {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, tc.exp, res)
			}
		})
	}
}
