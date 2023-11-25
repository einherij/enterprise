package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	Region          string `mapstructure:"region"`
	// Used for minio as local aws
	Local    bool   `mapstructure:"local"`
	Endpoint string `mapstructure:"endpoint"`
}

type S3ClientAWS struct {
	Session *session.Session
	S3      *s3.S3
	Bucket  string
}

func NewS3ClientAWS(cfg AWSConfig) (*S3ClientAWS, error) {
	var s3ClientAWS = new(S3ClientAWS)
	s3ClientAWS.Bucket = cfg.Bucket
	sessionConfig := &aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}
	if cfg.Local {
		sessionConfig.Endpoint = aws.String(cfg.Endpoint)
		sessionConfig.S3ForcePathStyle = aws.Bool(true)
		sessionConfig.DisableSSL = aws.Bool(true)
	}
	sess, err := session.NewSession(sessionConfig)
	if err != nil {
		return nil, err
	}
	s3ClientAWS.Session = sess
	s3ClientAWS.S3 = s3.New(sess)
	return s3ClientAWS, nil
}
