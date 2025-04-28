package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Logging `mapstructure:"logging"`
}

type Logging struct {
	Level     string `mapstructure:"level"`
	Formatter string `mapstructure:"formatter"`
	Handler   string `mapstructure:"handler"`
	Path      string `mapstructure:"path"`
}

type EnvConfig struct {
	APIHost string `env:"API_HOST" envDefault:"localhost"`
	APIPort int    `env:"API_PORT" envDefault:"3000"`

	DBHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	DBPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	DBUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	DBPassword string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"POSTGRES_DB" envDefault:"postgres"`

	MigrationsPath string `env:"POSTGRES_MIGRATIONS_PATH" envDefault:"./migrations"`
	MigrationsDSN  string `env:"POSTGRES_MIGRATIONS_DSN" envDefault:"host=$(DB_HOST) port=$(DB_PORT) dbname=$(DB_NAME) user=$(DB_USER) password=$(DB_PASSWORD) sslmode=disable"`

	RedisHost string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisDB   int    `env:"REDIS_DB" envDefault:"0"`
}

// LoadConfig reads configuration from yml file
func LoadConfig(path string) (*Config, error) {

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err.Error())
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshal config: %v", err.Error())
	}
	return config, nil
}

// LoadEnvConfig reads configuration from env file
func LoadEnvConfig() (*EnvConfig, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err.Error())
	}

	envCfg := new(EnvConfig)
	if err := env.Parse(envCfg); err != nil {
		return nil, fmt.Errorf("unable to parse ennvironment variables: %v", err.Error())
	}
	return envCfg, nil
}
