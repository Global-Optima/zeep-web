package config

type FilteringConfig struct {
	DefaultPage     int `mapstructure:"DEFAULT_PAGE" default:"1"`
	DefaultPageSize int `mapstructure:"DEFAULT_PAGE_SIZE" default:"10"`
	MaxPageSize     int `mapstructure:"MAX_PAGE_SIZE" default:"100"`
}
