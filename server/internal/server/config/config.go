package config



type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func NewConfig(host, port string) *Config {
	return &Config{
		Host: host,
		Port: port,
	}
}