package config

type TTLConfig struct {
	Hot  int `mapstructure:"TTL_HOT"`
	Warm int `mapstructure:"TTL_WARM"`
	Cool int `mapstructure:"TTL_COLD"`
}
