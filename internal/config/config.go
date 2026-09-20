package config

import "os"

type Config struct {
	StorageType string
	DatabaseURL string
}

func Load() Config {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "memory"
	}

	return Config{
		StorageType: storageType,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
