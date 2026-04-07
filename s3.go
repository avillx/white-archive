package whitearchive

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client     *minio.Client
	bucketName string
}

func NewStorageClient(endpoint, accessKey, secretKey, bucket string) (*StorageClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return &StorageClient{client: client, bucketName: bucket}, nil
}

func (s *StorageClient) Upload(ctx context.Context, name string, data []byte) error {
	_, err := s.client.PutObject(ctx, s.bucketName, name, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (s *StorageClient) Download(ctx context.Context, name string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)

	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return []byte{}, ErrNotFound
		}
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrEmptyFile
	}

	return data, nil
}

// func (s *StorageClient) Delete(ctx context.Context, name string) error {
// 	return s.client.RemoveObject(ctx, s.bucketName, name, minio.RemoveObjectOptions{})
// }
