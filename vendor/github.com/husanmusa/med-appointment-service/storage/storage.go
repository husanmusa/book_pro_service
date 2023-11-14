package storage

import (
	"context"
	pb "github.com/husanmusa/med-appointment-service/genproto/appointment_service"
)

type StorageI interface {
	Appointment() AppointmentServiceI
	Doctor() DoctorI
	Patient() PatientI
	CloseDB()
}

type AppointmentServiceI interface {
	Create(ctx context.Context, req *pb.Appointment) (*pb.Appointment, error)
	Get(ctx context.Context, req *pb.AppointmentId) (*pb.Appointment, error)
	Cancel(ctx context.Context, req *pb.AppointmentId) error
}

type DoctorI interface {
	Create(ctx context.Context, in *pb.Doctor) error
	GetList(ctx context.Context, in *pb.GetListDoctorsRequest) (*pb.GetListDoctorsResponse, error)
	Get(ctx context.Context, in *pb.DoctorId) (*pb.Doctor, error)
	Update(ctx context.Context, in *pb.Doctor) error
	Delete(ctx context.Context, in *pb.DoctorId) error
}

type PatientI interface {
	Create(ctx context.Context, in *pb.Patient) error
	GetList(ctx context.Context, in *pb.GetListPatientsRequest) (*pb.GetListPatientsResponse, error)
	Get(ctx context.Context, in *pb.PatientId) (*pb.Patient, error)
	Update(ctx context.Context, in *pb.Patient) error
	Delete(ctx context.Context, in *pb.PatientId) error
}
