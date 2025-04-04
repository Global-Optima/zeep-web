package config

type S3Config struct {
	AccessKey        string `mapstructure:"S3_ACCESS_KEY" validate:"required"`
	SecretKey        string `mapstructure:"S3_SECRET_KEY" validate:"required"`
	AccessEndpoint   string `mapstructure:"S3_ACCESS_ENDPOINT" validate:"required"`
	ResponseEndpoint string `mapstructure:"S3_RESPONSE_ENDPOINT" validate:"required"`
	BucketName       string `mapstructure:"S3_BUCKET_NAME" validate:"required"`
}
