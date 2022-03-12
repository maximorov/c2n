package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	TelegramToken string
}

var Cnf Config

// Trying to initialize env variables from file and checks if passed ENV type exist
func InitEnv(filesDir string) {
	envFilePath, _ := filepath.Abs(filesDir + ".env")

	if err := godotenv.Load(envFilePath); err != nil {
		zap.S().Warn("Error in loading {.env} files")
	}
}

func initConfig() {
	Cnf = Config{
		Env:           os.Getenv("ENV"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}
}

func initLogger() {
	var config zap.Config
	if Cnf.Env == `dev` {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()

	zap.ReplaceGlobals(logger)
}
