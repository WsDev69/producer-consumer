package common

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/WsDev69/producer-consumer/pkg/build"
)

func ShowVersion() {
	versionFlag := flag.Bool("version", false, "prints the build version")
	flag.Parse()

	// If the version flag is passed, print the version and exit
	if *versionFlag {
		fmt.Printf("Version: %s\n", build.Version)
		os.Exit(0)
	}
}

func Logger(logLevel, output string) *slog.Logger {
	slogLevel := logLevelToSlog(logLevel)
	slogHandler := getHandler(output, slogLevel)

	logger := slog.New(slogHandler)

	slog.SetDefault(logger)

	return logger
}

func getHandler(output string, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	switch output {
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		return slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.NewJSONHandler(os.Stdout, opts)
}

func logLevelToSlog(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelDebug
}
