package service

import (
	"context"
	pb "github.com/husanmusa/med-appointment-service/genproto/appointment_service"
	"github.com/husanmusa/med-appointment-service/pkg/helper"
	"github.com/husanmusa/med-appointment-service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type doctorService struct {
	logger  logger.LoggerI
	storage storage.StorageI
	pb.UnimplementedDoctorServiceServer
}

func NewDoctorService(log logger.LoggerI, store storage.StorageI) *doctorService {
	return &doctorService{
		logger:  log,
		storage: store,
	}
}

func (s *doctorService) Create(ctx context.Context, req *pb.Doctor) (*emptypb.Empty, error) {
	s.logger.Info("---Create Destrict--->", logger.Any("req", req))
	err := s.storage.Doctor().Create(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error creating doctor", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *doctorService) GetList(ctx context.Context, req *pb.GetListDoctorsRequest) (*pb.GetListDoctorsResponse, error) {
	s.logger.Info("---Get List Doctor--->", logger.Any("req", req))
	resp, err := s.storage.Doctor().GetList(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting list doctor", req, codes.Internal)
	}
	return resp, nil
}

func (s *doctorService) Get(ctx context.Context, req *pb.DoctorId) (*pb.Doctor, error) {
	s.logger.Info("---Get Doctor--->", logger.Any("req", req))
	resp, err := s.storage.Doctor().Get(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting doctor", req, codes.Internal)
	}
	return resp, nil
}

func (s *doctorService) Update(ctx context.Context, req *pb.Doctor) (*emptypb.Empty, error) {
	s.logger.Info("---Update Doctors--->", logger.Any("req", req))
	err := s.storage.Doctor().Update(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error Updating doctor", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *doctorService) Delete(ctx context.Context, req *pb.DoctorId) (*emptypb.Empty, error) {
	s.logger.Info("---Delete Doctors--->", logger.Any("req", req))
	err := s.storage.Doctor().Delete(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error deleting doctor", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}
