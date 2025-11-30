package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string     `yaml:"env"`
	JWTSecret  []byte     `yaml:"jwt_secret"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Database   Database   `yaml:"database"`
	CDEK       CDEK       `yaml:"cdek"`
	Dellin     Dellin     `yaml:"dellin"`
}

type HTTPServer struct {
	Port         string        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	DSN             string        `yaml:"dsn"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns" default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" default:"1h"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" default:"30m"`
}

type CDEK struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Login        string `yaml:"login"`
	Password     string `yaml:"password"`
}

type Dellin struct {
	AppKey string `yaml:"app_key"`
}

var configPath string = "./config/config.yaml"

func MustLoadConfig() *Config {

	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetTypeByDefaultValue(true)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	config := &Config{
		Env:       v.GetString("env"),
		JWTSecret: []byte(v.GetString("jwt_secret")),
		HTTPServer: HTTPServer{
			Port:         v.GetString("http_server.port"),
			WriteTimeout: v.GetDuration("http_server.write_timeout"),
			ReadTimeout:  v.GetDuration("http_server.read_timeout"),
			IdleTimeout:  v.GetDuration("http_server.idle_timeout"),
		},
		Database: Database{
			DSN:             v.GetString("database.dsn"),
			Host:            v.GetString("database.host"),
			Port:            v.GetInt("database.port"),
			User:            v.GetString("database.user"),
			Password:        v.GetString("database.password"),
			Name:            v.GetString("database.name"),
			SSLMode:         v.GetString("database.sslmode"),
			MaxOpenConns:    v.GetInt("database.max_open_conns"),
			MaxIdleConns:    v.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: v.GetDuration("database.conn_max_lifetime"),
			ConnMaxIdleTime: v.GetDuration("database.conn_max_idle_time"),
		},
		CDEK: CDEK{
			ClientID:     v.GetString("cdek.client_id"),
			ClientSecret: v.GetString("cdek.client_secret"),
			Login:        v.GetString("cdek.login"),
			Password:     v.GetString("cdek.password"),
		},
		Dellin: Dellin{
			AppKey: v.GetString("dellin.app_key"),
		},
	}

	return config
}

func (db *Database) GetDSN() string {
	if db.DSN == "" {
		return db.DSN
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.Name,
		db.SSLMode,
	)
}
