package aws

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

type S3Client struct {
	presign  *s3.PresignClient
	bucket   string
	mockMode bool
}

func NewS3Client(cfg awssdk.Config, bucket string, mockMode bool) *S3Client {
	var pc *s3.PresignClient
	if !mockMode {
		pc = s3.NewPresignClient(s3.NewFromConfig(cfg))
	}
	return &S3Client{presign: pc, bucket: bucket, mockMode: mockMode}
}

func (s *S3Client) PresignPutURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if s.mockMode {
		logger.WithContext(ctx).Info("[S3 MOCK] presign put URL",
			zap.String("bucket", s.bucket),
			zap.String("key", key),
		)
		return fmt.Sprintf("https://mock-s3.local/%s/%s", s.bucket, key), nil
	}

	req, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: awssdk.String(s.bucket),
		Key:    awssdk.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("s3: presign put %s: %w", key, err)
	}
	return req.URL, nil
}
