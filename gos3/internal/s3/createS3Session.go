package s3

import (
	"gos3/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func createS3Session(cfg config.S3Config) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
	})
}
