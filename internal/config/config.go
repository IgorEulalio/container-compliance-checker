package config

type Config struct {
	LogLevel string  `mapstructure:"log_level"`
	Checks   []Check `mapstructure:"checks"`
}

type Check struct {
	Type   string                 `mapstructure:"type"`
	Config map[string]interface{} `mapstructure:"config"`
}
