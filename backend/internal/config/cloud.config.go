package config

type S3Config struct {
	AccessKey  string `mapstructure:"PSKZ_ACCESS_KEY"`
	SecretKey  string `mapstructure:"PSKZ_SECRET_KEY"`
	Endpoint   string `mapstructure:"PSKZ_ENDPOINT"`
	BucketName string `mapstructure:"PSKZ_BUCKETNAME"`
}
