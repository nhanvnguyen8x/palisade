package config

import "fmt"

func (c *Config) Validate() error {
	if c.Storage.Endpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT is required")
	}

	if c.Storage.AccessKey == "" {
		return fmt.Errorf("MINIO_ACCESS_KEY is required")
	}

	if c.Storage.SecretKey == "" {
		return fmt.Errorf("MINIO_SECRET_KEY is required")
	}

	if c.Storage.Bucket == "" {
		return fmt.Errorf("MINIO_BUCKET is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}
