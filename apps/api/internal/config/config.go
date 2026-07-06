package config

type Config struct {
	HTTP    HTTP
	Storage Storage
}

type HTTP struct {
	Port int `env:"HTTP_PORT" envDefault:"8080"`
}

type Storage struct {
	Endpoint  string `env:"MINIO_ENDPOINT"`
	AccessKey string `env:"MINIO_ACCESS_KEY"`
	SecretKey string `env:"MINIO_SECRET_KEY"`
	Bucket    string `env:"MINIO_BUCKET"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
}
