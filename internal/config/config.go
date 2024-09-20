package config

import "time"

type (
	Config struct {
		LogLevel    string `default:"debug"`
		OutPut      string `default:"json"`
		MessageRate int    `default:"100"`
		Postgres    Postgres
		GRPC        GRPC
		Prometheus  Prometheus
	}

	Consumer struct {
		Config
	}

	Producer struct {
		Config
		NumWorker  int   `default:"5"`
		MaxBacklog int64 `default:"10"`
	}

	Postgres struct {
		Host            string        `default:"127.0.0.1"`
		Port            int           `default:"5432"`
		Password        string        `default:"12345"`
		User            string        `default:"postgres"`
		DB              string        `default:"testdb"`
		ConnMaxLifeTime time.Duration `default:"5m"`
		MaxIdleConn     int           `default:"10"`
		MaxOpenConn     int           `default:"10"`
	}

	GRPC struct {
		Host     string `default:"127.0.0.1"`
		Port     int    `default:"8080"`
		TLS      bool   `default:"false"`
		CertFile string
		KeyFile  string
	}

	Prometheus struct {
		Port int `default:"2112"`
	}
)
