package config

import "time"

var C Config = Config{
	Redis: Redis{
		Host:     "redis:6379",
		Database: 0,
	},
}

// GetConfig -
func GetConfig() Config {
	return C
}

// Config -
type Config struct {
	Version string

	// sever port
	Port string `json:"port"`

	// Mod - dev / pro
	Mod string `json:"mod"`

	// Logger
	Logger string `json:"logger"`

	// Redis
	Redis Redis `json:"redis"`
}

// Redis - Redis 資料庫配置
type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Password string `json:"password" yaml:"password"`
	Database int    `json:"database" yaml:"database"`
	TTL      int64  `json:"ttl" yaml:"ttl"`
	DatabaseOption
}

// DatabaseOption - 資料庫額外參數
type DatabaseOption struct {
	SSL               bool          `json:"ssl" yaml:"ssl"`
	MaxPoolSize       uint64        `json:"maxPoolSize" yaml:"maxPoolSize"`
	MinPoolSize       uint64        `json:"minPoolSize" yaml:"minPoolSize"`
	MaxRetries        int64         `json:"maxRetries" yaml:"maxReties"`
	MaxIdelConns      int64         `json:"maxIdelConns" yaml:"maxIdelConns"`
	MinIdelConns      uint64        `json:"minIdelConns" yaml:"minIdelConns"`
	MaxConns          int64         `json:"maxConns" yaml:"maxConns"`
	MinConns          int64         `json:"minConns" yaml:"minConns"`
	MaxConnIdleTime   time.Duration `json:"maxConnIdleTime" yaml:"maxConnIdleTime"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval" yaml:"heartbeatInterval"`
	Direct            bool          `json:"direct" yaml:"direct"`
}
