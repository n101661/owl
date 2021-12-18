package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	flags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Flag struct {
	Development bool   `long:"dev"`
	Directory   string `long:"dir" description:"the directory to load configurations" required:"true"`
}

func main() {
	var flag Flag
	_, err := flags.Parse(&flag)
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	logger, err := newLogger(flag.Development)
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	crons := newCron()

	err = filepath.Walk(flag.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !isYAMLExtension(path) {
			return nil
		}

		// load config
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		return crons.AddFromFile(f)
	})
	if err != nil {
		logger.Fatal("failed to load configurations", zap.Error(err))
	}

	if err := crons.StartAll(); err != nil {
		logger.Fatal("failed to start schedules", zap.Error(err))
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)
	<-sig

	crons.Shutdown()
}

func newLogger(development bool) (*zap.Logger, error) {
	var cfg zap.Config
	if development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
		cfg.DisableCaller = true
	}
	cfg.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(time.RFC3339Nano))
	}
	return cfg.Build()
}

func isYAMLExtension(path string) bool {
	return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
}
