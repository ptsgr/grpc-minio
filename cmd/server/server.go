package main

import (
	"log"
	"net"
	"time"

	"github.com/ptsgr/grpc-minio/internal/filestorage"
	"github.com/ptsgr/grpc-minio/internal/minfs"
	pb "github.com/ptsgr/grpc-minio/internal/server_grpc"
	"github.com/ptsgr/grpc-minio/internal/service"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	grpcPort = "grpc.port"

	minioURIKey        = "minio.uri"
	minioAccessKey     = "minio.access-key"
	minioSecretKey     = "minio.access-secret"
	minioBucketNameKey = "minio.bucket-name"
	minioUseSSLKey     = "minio.use-ssl"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}
	fileStorageBucket := viper.GetString(minioBucketNameKey)

	fileStorageClient, err := minfs.NewClient(minfs.Config{
		AccessKeyID:     viper.GetString(minioAccessKey),
		Endpoint:        viper.GetString(minioURIKey),
		SecretAccessKey: viper.GetString(minioSecretKey),
		UseSSL:          viper.GetBool(minioUseSSLKey),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = fileStorageClient.MakeBucket(fileStorageBucket)
	if err != nil {
		log.Fatal(err)
	}

	presignedURLTTL := time.Minute * 5
	fileStorage := filestorage.NewRemote(fileStorageClient, fileStorageBucket, presignedURLTTL)
	svc := service.NewService(fileStorage)
	listener, err := net.Listen("tcp", ":"+viper.GetString(grpcPort))

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterServerServer(grpcServer, svc)
	log.Println("Starting grpc server on port: " + viper.GetString(grpcPort))
	grpcServer.Serve(listener)
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("local")
	return viper.ReadInConfig()
}
