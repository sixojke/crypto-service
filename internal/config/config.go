package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServer
	Postgres   Postgres `mapstructure:"postgres"`
	Logger     Logger
}

type HTTPServer struct {
	Port               string        `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	MaxHeaderMegabytes int           `mapstructure:"max_header_megabytes"`
}

type Postgres struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	MigrationsPath  string `mapstructure:"migrations_path"`
}

type Logger struct {
	LogLevel int `mapstructure:"log_level"`
}

func Init(configPaths []string, envFile string) (*Config, error) {
	var cfg Config

	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("loading env file: %w", err)
		}
	}

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	if err := godotenv.Load(envFile); err != nil {
		return nil, fmt.Errorf("failed to read env file: %s: %s", envFile, err)
	}

	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")

	return &cfg, nil
}
