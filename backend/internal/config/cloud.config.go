package config

type S3Config struct {
	AccessKey  string `mapstructure:"S3_ACCESS_KEY"`
	SecretKey  string `mapstructure:"S3_SECRET_KEY"`
	Endpoint   string `mapstructure:"S3_ENDPOINT"`
	BucketName string `mapstructure:"S3_BUCKET_NAME"`
}
