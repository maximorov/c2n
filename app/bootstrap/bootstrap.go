package bootstrap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	TelegramToken string
	DB            ConfigDb
}

type ConfigDb struct {
	Host                 string        `env:"HOST" validate:"min=3"`
	User                 string        `env:"USER" validate:"min=3"`
	Password             string        `env:"PASSWORD" validate:"min=3"`
	DBName               string        `env:"NAME" validate:"min=3"`
	Port                 uint          `env:"PORT" validate:"min=1"`
	MaxOpenCons          uint          `env:"MAX_OPEN_CONNS" validate:"min=1,bound-cons"`
	MaxIdleCons          uint          `env:"MAX_IDLE_CONNS" validate:"min=0,bound-idle"`
	MaxLifetimeMinutes   time.Duration `env:"MAX_LIFETIME_MINUTES" validate:"min=1m"`
	MaxIdleMinutes       time.Duration `env:"MAX_IDLE_MINUTES" validate:"min=0m"`
	Log                  uint          `env:"LOG" validate:"lte=6"`
	PreferSimpleProtocol bool          `env:"PREFER_SIMPLE_PROTOCOL" envDefault:"true"`
	SSLMode              string        `env:"SSL_MODE" envDefault:"disable" validate:"oneof=disable allow default require verify-ca verify-full"`
	Timezone             string        `env:"TZ" envDefault:"UTC" validate:"timezone"`
}

func (d ConfigDb) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d pool_min_conns=%d "+
			"prefer_simple_protocol=%t sslmode=%s timezone=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s ",
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.DBName,
		d.MaxOpenCons,
		d.MaxIdleCons,
		d.PreferSimpleProtocol,
		d.SSLMode,
		d.Timezone,
		d.MaxLifetimeMinutes,
		d.MaxIdleMinutes,
	)
}

var Cnf Config

// Trying to initialize env variables from file and checks if passed ENV type exist
func InitEnv(filesDir string) {
	envFilePath, _ := filepath.Abs(filesDir + ".env")

	if err := godotenv.Load(envFilePath); err != nil {
		zap.S().Warn("Error in loading {.env} files")
	}
}

func InitConfig() {
	Cnf = Config{
		Env:           os.Getenv("ENV"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DB:            ConfigDb{},
	}
}

func InitLogger() {
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
