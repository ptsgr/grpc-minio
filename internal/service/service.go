package service

import (
	"context"
	"log"

	pb "github.com/ptsgr/grpc-minio/internal/server_grpc"
)

type service struct {
	fileStorage FileStorage
}

func NewService(
	storage FileStorage,
) *service {
	return &service{
		fileStorage: storage,
	}
}

func (s *service) Put(c context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	fileName, err := s.fileStorage.Put(c, request.Filename, request.FileData)
	if err != nil {
		log.Println(err)
		return &pb.PutResponse{
			Message: "faild to put file: " + err.Error(),
			Code:    pb.StatusCode_Failed,
		}, err
	}
	return &pb.PutResponse{
		Message:  "file uploaded successfully",
		Code:     pb.StatusCode_Ok,
		Filename: fileName,
	}, nil
}

func (s *service) Get(c context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	fileData, err := s.fileStorage.Get(c, request.Filename)
	if err != nil {
		log.Println(err)
		return &pb.GetResponse{
			Message: "faild to get file: " + err.Error(),
			Code:    pb.StatusCode_Failed,
		}, err
	}
	return &pb.GetResponse{
		Message:  "successfully get file",
		Code:     pb.StatusCode_Ok,
		FileData: fileData,
	}, nil
}

func (s *service) Delete(c context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.fileStorage.Remove(c, request.Filename)
	if err != nil {
		log.Println(err)
		return &pb.DeleteResponse{
			Message: "faild to delete file: " + err.Error(),
			Code:    pb.StatusCode_Failed,
		}, err
	}
	return &pb.DeleteResponse{
		Message: "successfully delete file",
		Code:    pb.StatusCode_Ok,
	}, nil
}
