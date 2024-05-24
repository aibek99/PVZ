package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

// Config is
type Config struct {
	Server struct {
		GRPCPort    string `json:"GRPCPort"`
		HTTPPort    string `json:"HTTPPort"`
		MetricsPort string `json:"MetricsPort"`
	}
	Postgres Postgres `json:"Postgres"`
	Auth     struct {
		User     string `json:"AuthUser"`
		Password string `json:"Password"`
	} `json:"Auth"`
	ShutdownTime time.Duration `json:"ShutdownTime"`
	TLS          struct {
		CertPath string `json:"CertKeyPath"`
		KeyPath  string `json:"KeyPath"`
		CACrt    string `json:"CACrt"`
	} `json:"TLS"`
	Kafka         Kafka         `json:"Kafka"`
	Redis         Redis         `json:"Redis"`
	InMemoryCache InMemoryCache `json:"InMemoryCache"`
	CacheType     string        `json:"CacheType"`
}

// Postgres is
type Postgres struct {
	Host     string `json:"Host" validate:"required"`
	Port     string `json:"Port" validate:"required"`
	User     string `json:"User" validate:"required"`
	Password string `json:"Password" validate:"required"`
	DBName   string `json:"DBName" validate:"required"`
}

// Kafka is
type Kafka struct {
	Topic   string   `json:"topic" validate:"required"`
	Brokers []string `json:"brokers" validate:"required"`
}

// Redis is
type Redis struct {
	Address  string `json:"address" validate:"required"`
	Password string `json:"password"`
}

// InMemoryCache is
type InMemoryCache struct {
	CleanTime float64 `json:"CleanTime"`
}

// LoadConfig is
func LoadConfig(configPath string) (*Config, error) {
	// #nosec G304
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	defer func() {
		err = jsonFile.Close()
		if err != nil {
			log.Printf("[config][LoadConfig] jsonFile.Close: %v\n", err)
		}
	}()

	var c Config

	err = json.NewDecoder(jsonFile).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	err = validator.New().Struct(c)
	if err != nil {
		return nil, fmt.Errorf("validator.New.Struct: %w", err)
	}

	return &c, nil
}
