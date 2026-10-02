package config

import "os"

type Config struct {
	MySQLDSN string
	JWTKey   string
	Port     string
}

func Load() *Config {
	return &Config{
		MySQLDSN: getEnv("MYSQL_DSN", "root:12345678@tcp(127.0.0.1:3306)/recycle_trade?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTKey:   getEnv("JWT_KEY", "recycle-trade-secret-key-2024"),
		Port:     getEnv("PORT", "9000"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
