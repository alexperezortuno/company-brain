package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type BrainConfig struct {
	Instance struct {
		Name     string `yaml:"name"`
		Language string `yaml:"language"`
	} `yaml:"instance"`
	Modules map[string]ModuleConfig `yaml:"modules"`
}

type ModuleConfig struct {
	Enabled bool `yaml:"enabled"`
}

type RuntimeConfig struct {
	HTTPAddress       string
	ConfigPath        string
	DependencyTimeout time.Duration
	PostgresAddress   string
	RedisAddress      string
	RedisPassword     string
	QdrantURL         string
	MinIOHealthURL    string
	WorkerHealthURL   string
	Brain             BrainConfig
}

func Load() (RuntimeConfig, error) {
	cfg := RuntimeConfig{
		HTTPAddress:       env("BRAIN_HTTP_ADDRESS", ":8080"),
		ConfigPath:        env("BRAIN_CONFIG_PATH", "/app/config/brain.yaml"),
		PostgresAddress:   env("POSTGRES_ADDRESS", "postgres:5432"),
		RedisAddress:      env("REDIS_ADDRESS", "redis:6379"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		QdrantURL:         env("QDRANT_URL", "http://qdrant:6333/"),
		MinIOHealthURL:    env("MINIO_HEALTH_URL", "http://minio:9000/minio/health/live"),
		WorkerHealthURL:   env("WORKER_HEALTH_URL", "http://knowledge-worker:8090/health"),
		DependencyTimeout: durationEnv("DEPENDENCY_TIMEOUT", 2*time.Second),
	}

	raw, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("read brain config: %w", err)
	}
	if err := yaml.Unmarshal(raw, &cfg.Brain); err != nil {
		return RuntimeConfig{}, fmt.Errorf("parse brain config: %w", err)
	}
	if cfg.Brain.Instance.Name == "" {
		return RuntimeConfig{}, errors.New("instance.name is required")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	if value, err := time.ParseDuration(raw); err == nil {
		return value
	}
	if seconds, err := strconv.Atoi(raw); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return fallback
}
