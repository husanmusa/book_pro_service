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

type patientService struct {
	logger  logger.LoggerI
	storage storage.StorageI
	pb.UnimplementedPatientServiceServer
}

func NewPatientService(log logger.LoggerI, store storage.StorageI) *patientService {
	return &patientService{
		logger:  log,
		storage: store,
	}
}

func (s *patientService) Create(ctx context.Context, req *pb.Patient) (*emptypb.Empty, error) {
	s.logger.Info("---Create Destrict--->", logger.Any("req", req))
	err := s.storage.Patient().Create(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error creating patient", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *patientService) GetList(ctx context.Context, req *pb.GetListPatientsRequest) (*pb.GetListPatientsResponse, error) {
	s.logger.Info("---Get List Patient--->", logger.Any("req", req))
	resp, err := s.storage.Patient().GetList(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting list patient", req, codes.Internal)
	}
	return resp, nil
}

func (s *patientService) Get(ctx context.Context, req *pb.PatientId) (*pb.Patient, error) {
	s.logger.Info("---Get Patient--->", logger.Any("req", req))
	resp, err := s.storage.Patient().Get(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting patient", req, codes.Internal)
	}
	return resp, nil
}

func (s *patientService) Update(ctx context.Context, req *pb.Patient) (*emptypb.Empty, error) {
	s.logger.Info("---Update Patients--->", logger.Any("req", req))
	err := s.storage.Patient().Update(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error Updating patient", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *patientService) Delete(ctx context.Context, req *pb.PatientId) (*emptypb.Empty, error) {
	s.logger.Info("---Delete Patients--->", logger.Any("req", req))
	err := s.storage.Patient().Delete(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error deleting patient", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}
