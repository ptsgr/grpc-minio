package filestorage

import (
	"bytes"
	"context"
	"time"

	"github.com/ptsgr/grpc-minio/internal/docid"
	"github.com/ptsgr/grpc-minio/internal/service"
)

type Remote struct {
	client     service.FileStorageClient
	bucketName string
	ttl        time.Duration
}

func NewRemote(client service.FileStorageClient, bucketName string, ttl time.Duration) *Remote {
	return &Remote{client: client, bucketName: bucketName, ttl: ttl}
}

func (r *Remote) Put(ctx context.Context, fileName string, fileData []byte) (string, error) {
	fileName = docid.New() + "/" + fileName
	return fileName, r.client.Write(r.bucketName, fileName, bytes.NewReader(fileData))
}

func (r *Remote) Get(ctx context.Context, fileName string) ([]byte, error) {
	return r.client.ReadBytes(r.bucketName, fileName)
}

func (r *Remote) Remove(ctx context.Context, fileName string) error {
	return r.client.Remove(r.bucketName, fileName)
}
