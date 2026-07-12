package config

import "strconv"

type Config struct {
	HTTP     HTTP
	Database Database
	Storage  Storage
}

type HTTP struct {
	Port int `env:"HTTP_PORT" envDefault:"8080"`
}

type Database struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"palisade"`
	Password string `env:"DB_PASSWORD" envDefault:"palisade"`
	Name     string `env:"DB_NAME" envDefault:"palisade"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

func (d Database) URL() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + strconv.Itoa(d.Port) + "/" + d.Name + "?sslmode=" + d.SSLMode
}

type Storage struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"admin"`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:"password123"`
	Bucket    string `env:"MINIO_BUCKET" envDefault:"palisade"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
}
