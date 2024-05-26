package config

import (
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"
)

var (
	once   sync.Once
	config Config
)

// Config struct holds the configuration for the application.
type Config struct {
	ConnDB     string
	LogLevelDB string
	Log        struct {
		Level zerolog.Level
		JSON  bool
	}
	JWT struct {
		Secret   []byte
		Duration time.Duration
	}
	SMTP struct {
		Host     string
		Port     uint16
		Username string
		Password string
	}
	Cache struct {
		ResetPass time.Duration
		Activate  time.Duration
	}
	Domain      string
	ServerProps struct {
		Port                 string
		ReadTimeout          time.Duration
		WriteTimeout         time.Duration
		IdleTimeout          time.Duration
		CancelContextTimeout time.Duration
	}
	SnowflakeNode int64
}

type envs struct {
	DBHost               string        `env:"DB_HOST,required,notEmpty"`
	DBPort               uint16        `env:"DB_PORT,required,notEmpty"`
	DBUser               string        `env:"DB_USER,required,notEmpty"`
	DBPassword           string        `env:"DB_PASSWORD,required,notEmpty"`
	DBName               string        `env:"DB_NAME,required,notEmpty"`
	DBSSLMode            string        `env:"DB_SSL_MODE" envDefault:"disable"`
	DBLogLevel           string        `env:"DB_LOG_LEVEL" envDefault:"warn"`
	LogLevel             string        `env:"LOG_LEVEL" envDefault:"info"`
	LogJSON              bool          `env:"LOG_JSON" envDefault:"true"`
	JWTSecret            string        `env:"JWT_SECRET,required,notEmpty"`
	JWTDuration          time.Duration `env:"JWT_DURATION" envDefault:"48h"`
	SMTPHost             string        `env:"SMTP_HOST,required,notEmpty"`
	SMTPPort             uint16        `env:"SMTP_PORT,required,notEmpty"`
	SMTPUser             string        `env:"SMTP_USER,required,notEmpty"`
	SMTPPass             string        `env:"SMTP_PASS,required,notEmpty"`
	DOMAIN               string        `env:"DOMAIN,required,notEmpty"`
	ServerPort           uint16        `env:"SERVER_PORT,required,notEmpty"`
	ReadTimeout          time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout         time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout          time.Duration `env:"IDLE_TIMEOUT" envDefault:"15s"`
	CancelContextTimeout time.Duration `env:"CANCEL_CONTEXT_TIMEOUT" envDefault:"30s"`
	SnowflakeNode        int64         `env:"SNOWFLAKE_NODE" envDefault:"1"`
}

// MustGet loads the configuration from environment variables.
func MustGet() *Config {
	once.Do(func() {
		slog.Info("loading config from environment variables")

		var e envs
		if err := env.Parse(&e); err != nil {
			log.Fatal(err)
		}

		e.jwt()
		e.connDb()
		e.logLevelDb()
		e.logger()
		e.sendMail()
		e.server()
		e.domain()
		e.snowflake()
	})

	return &config
}

func (e *envs) server() {
	config.ServerProps.Port = fmt.Sprintf(":%d", e.ServerPort)
	config.ServerProps.ReadTimeout = e.ReadTimeout
	config.ServerProps.WriteTimeout = e.WriteTimeout
	config.ServerProps.IdleTimeout = e.IdleTimeout
	config.ServerProps.CancelContextTimeout = e.CancelContextTimeout
}

func (e *envs) sendMail() {
	config.SMTP.Host = e.SMTPHost
	config.SMTP.Port = e.SMTPPort
	config.SMTP.Username = e.SMTPUser
	config.SMTP.Password = e.SMTPPass
}

func (e *envs) jwt() {
	config.JWT.Secret = []byte(e.JWTSecret)
	config.JWT.Duration = e.JWTDuration
}

func (e *envs) logger() {
	level, err := zerolog.ParseLevel(e.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	config.Log.Level = level
	config.Log.JSON = e.LogJSON
}

func (e *envs) connDb() {
	config.ConnDB = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		e.DBUser,
		e.DBPassword,
		e.DBHost,
		e.DBPort,
		e.DBName,
		e.DBSSLMode,
	)
}

func (e *envs) domain() {
	config.Domain = e.DOMAIN
}

func (e *envs) logLevelDb() {
	config.LogLevelDB = e.DBLogLevel
}

func (e *envs) snowflake() {
	config.SnowflakeNode = e.SnowflakeNode
}
