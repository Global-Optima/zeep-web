package config

type S3Config struct {
	AccessKey  string `mapstructure:"S3_ACCESS_KEY" validate:"required"`
	SecretKey  string `mapstructure:"S3_SECRET_KEY" validate:"required"`
	Endpoint   string `mapstructure:"S3_ENDPOINT" validate:"required"`
	BucketName string `mapstructure:"S3_BUCKET_NAME" validate:"required"`
}
