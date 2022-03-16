package bootstrap

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type Config struct {
	Env           string
	TelegramToken string
	DB            ConfigDb
}

type ConfigDb struct {
	Host                 string `env:"HOST" validate:"min=3"`
	User                 string `env:"USER" validate:"min=3"`
	Password             string `env:"PASSWORD" validate:"min=3"`
	DBName               string `env:"NAME" validate:"min=3"`
	Port                 string `env:"PORT" validate:"min=1"`
	MaxOpenCons          string `env:"MAX_OPEN_CONNS" validate:"min=1,bound-cons"`
	MaxIdleCons          string `env:"MAX_IDLE_CONNS" validate:"min=0,bound-idle"`
	MaxLifetimeMinutes   string `env:"MAX_LIFETIME_MINUTES" validate:"min=1m"`
	MaxIdleMinutes       string `env:"MAX_IDLE_MINUTES" validate:"min=0m"`
	Log                  string `env:"LOG" validate:"lte=6"`
	PreferSimpleProtocol string `env:"PREFER_SIMPLE_PROTOCOL" envDefault:"true"`
	SSLMode              string `env:"SSL_MODE" envDefault:"disable" validate:"oneof=disable allow default require verify-ca verify-full"`
	Timezone             string `env:"TZ" envDefault:"UTC" validate:"timezone"`
}

func (d ConfigDb) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s", /*+
		"pool_max_conns=%s pool_min_conns=%s "+
		"prefer_simple_protocol=%s sslmode=%s timezone=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s "*/
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.DBName,
		//d.MaxOpenCons,
		//d.MaxIdleCons,
		//d.PreferSimpleProtocol,
		//d.SSLMode,
		//d.Timezone,
		//d.MaxLifetimeMinutes,
		//d.MaxIdleMinutes,
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
		Env:           os.Getenv(`ENV`),
		TelegramToken: os.Getenv(`TELEGRAM_TOKEN`),
		DB: ConfigDb{
			Host:                 os.Getenv(`DB_HOST`),
			User:                 os.Getenv(`DB_USER`),
			Password:             os.Getenv(`DB_PASSWORD`),
			DBName:               os.Getenv(`DB_NAME`),
			Port:                 os.Getenv(`DB_PORT`),
			MaxOpenCons:          os.Getenv(`DB_MAX_OPEN_CONNS`),
			MaxIdleCons:          os.Getenv(`DB_MAX_IDLE_CONNS`),
			MaxLifetimeMinutes:   os.Getenv(`DB_MAX_LIFETIME_MINUTES`),
			MaxIdleMinutes:       os.Getenv(`DB_MAX_IDLE_MINUTES`),
			Log:                  os.Getenv(`DB_LOG`),
			PreferSimpleProtocol: `true`,
			SSLMode:              os.Getenv(`DB_SSL_MODE`),
			Timezone:             `UTC`,
		},
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
