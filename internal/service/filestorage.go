package service

import (
	"context"
	"io"
)

type FileStorage interface {
	Put(ctx context.Context, fileName string, fileData []byte) (string, error)
	Get(ctx context.Context, fileName string) ([]byte, error)
	Remove(ctx context.Context, fileName string) error
}

type FileStorageClient interface {
	Write(bucketName, path string, r io.Reader) error
	ReadBytes(bucketName, path string) ([]byte, error)
	Remove(bucketName, path string) error
}
