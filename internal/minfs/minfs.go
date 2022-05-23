package minfs

import (
	"context"
	"io"
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Bucket is common reduced interface for Minio actions.
// path is full path to file inside the bucket including file name: "my/file/here.txt"
type Bucket interface {
	Read(bucketName, path string) (io.ReadSeekCloser, error)
	Write(bucketName, path string, r io.Reader) error
	Remove(bucketName string, path string) error
}

// Client represents an opened minio connection.
type Client struct {
	client *minio.Client
	cfg    *Config
}

// New creates a new Client object.
func New(endpoint, AccessKeyID, SecretAccessKey string) (*Client, error) {
	return NewClient(Config{
		Endpoint:        endpoint,
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: SecretAccessKey,
	})
}

// NewClient creates a new Client object using config.
func NewClient(cfg Config) (*Client, error) {
	creds := credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, "")
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:        creds,
		Secure:       cfg.UseSSL,
		Region:       cfg.Region,
		BucketLookup: minio.BucketLookupPath,
	})

	c := Client{
		client: client,
		cfg:    &cfg,
	}
	return &c, err
}

// MakeBucket creates new bucket and does not fail if bucket already exists.
func (c *Client) MakeBucket(bucketName string) error {
	exists, errBucketExists := c.client.BucketExists(context.Background(), bucketName)
	if errBucketExists == nil && exists {
		return nil
	}
	if errBucketExists != nil {
		return errBucketExists
	}

	err := c.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes reads data for an object.
func (c *Client) ReadBytes(bucketName, path string) ([]byte, error) {
	ff, err := c.client.GetObject(context.Background(), bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer ff.Close()
	data, err := ioutil.ReadAll(ff)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Write data from Reader.
// It creates the file if it doesn't exist.
func (c *Client) Write(bucketName, path string, r io.Reader) error {
	return c.writeFrom(context.Background(), bucketName, path, r, -1)
}

// Remove removes a file for the given name.
func (c *Client) Remove(bucketName string, path string) error {
	path = cleanPath(path)
	return c.client.RemoveObject(context.Background(), bucketName, path, minio.RemoveObjectOptions{})
}

func (c *Client) writeFrom(ctx context.Context, bucketName, path string, r io.Reader, size int64) error {
	path = cleanPath(path)
	opts := minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	}

	_, err := c.client.PutObject(ctx, bucketName, path, r, size, opts)
	if err != nil {
		return err
	}
	return nil
}

func getContentType(path string) string {
	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func cleanPath(name string) string {
	return filepath.Clean(strings.Trim(name, "/"))
}
