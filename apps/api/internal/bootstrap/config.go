package bootstrap

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Port string
}

func LoadConfig() Config {
	return Config{
		HTTP: HTTPConfig{
			Port: ":8080",
		},
	}
}
