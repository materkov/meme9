package uploader

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Uploader struct {
	s3Client *s3.Client
}

func New(accessKeyID, secretAccessKey string) (*Uploader, error) {
	region := "ru-7"
	host := "https://s3.ru-7.storage.selcloud.ru"

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(host)
	})

	return &Uploader{
		s3Client: s3Client,
	}, nil
}

func (u *Uploader) Upload(ctx context.Context, file []byte) (string, error) {
	// Public bucket host example:
	// https://874a91e6-f73e-438f-b7a4-0aa63b0959f9.selstorage.ru/<object_key>
	publicBucketHost := "https://874a91e6-f73e-438f-b7a4-0aa63b0959f9.selstorage.ru"
	bucket := "meme-files"

	key := fmt.Sprintf("photos/%d.jpg", time.Now().UnixNano())

	_, err := u.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", publicBucketHost, key)
	return publicURL, nil
}
