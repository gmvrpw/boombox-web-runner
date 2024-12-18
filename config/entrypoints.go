package config

type EntrypointsConfig struct {
	HTTP HttpEntrypointConfig `yaml:"http"`
}

type HttpEntrypointConfig struct {
	Port int `yaml:"port"`
}
