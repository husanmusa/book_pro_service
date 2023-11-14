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

type appointmentService struct {
	logger  logger.LoggerI
	storage storage.StorageI
	pb.UnimplementedAppointmentServiceServer
}

func NewAppointmentService(log logger.LoggerI, store storage.StorageI) *appointmentService {
	return &appointmentService{
		logger:  log,
		storage: store,
	}
}

func (s *appointmentService) Create(ctx context.Context, req *pb.Appointment) (*pb.Appointment, error) {
	s.logger.Info("---Create Appointment--->", logger.Any("req", req))
	resp, err := s.storage.Appointment().Create(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error creating appointment", req, codes.Internal)
	}
	return resp, nil
}

func (s *appointmentService) Get(ctx context.Context, req *pb.AppointmentId) (*pb.Appointment, error) {
	s.logger.Info("---Get Appointment--->", logger.Any("req", req))
	resp, err := s.storage.Appointment().Get(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting appointment", req, codes.Internal)
	}
	return resp, nil
}

func (s *appointmentService) Cancel(ctx context.Context, req *pb.AppointmentId) (*emptypb.Empty, error) {
	s.logger.Info("---Cancel Appointment--->", logger.Any("req", req))
	err := s.storage.Appointment().Cancel(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error canceling appointment", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}
