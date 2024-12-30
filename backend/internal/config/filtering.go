package config

type FilteringConfig struct {
	DefaultPage          int    `mapstructure:"DEFAULT_PAGE"`
	DefaultPageSize      int    `mapstructure:"DEFAULT_PAGE_SIZE"`
	MaxPageSize          int    `mapstructure:"MAX_PAGE_SIZE"`
	DefaultSortParameter string `mapstructure:"DEFAULT_SORT_PARAMETER"`
	DefaultSortDirection string `mapstructure:"DEFAULT_SORT_DIRECTION"`
}
