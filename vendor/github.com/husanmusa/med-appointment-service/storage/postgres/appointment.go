package postgres

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/husanmusa/med-appointment-service/genproto/appointment_service"
	"github.com/husanmusa/med-appointment-service/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AppointmentRepo struct {
	Db *pgxpool.Pool
}

func NewAppointmentRepo(db *pgxpool.Pool) storage.AppointmentServiceI {
	return &AppointmentRepo{
		Db: db,
	}
}

func (a AppointmentRepo) Create(ctx context.Context, in *pb.Appointment) (*pb.Appointment, error) {
	err := a.Db.QueryRow(ctx, `insert into appointment (
                id,
                appoint_id,
                doctor_id,
                patient_id,
                appoint_date,
                appoint_time,
                status
                ) 
		values ($1, 
		coalesce((select appoint_id from appointment where appoint_date = to_char(current_timestamp, 'yyyy-mm-dd') order by appoint_id desc limit 1), 0)+1,
		 $2, $3, $4, $5, $6
		       ) returning appoint_id`,
		uuid.NewString(),
		in.DoctorId,
		in.PatientId,
		in.AppointDate,
		in.AppointTime,
		1,
	).Scan(&in.AppointId)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (a AppointmentRepo) Get(ctx context.Context, in *pb.AppointmentId) (*pb.Appointment, error) {
	var (
		appointment pb.Appointment
	)
	err := a.Db.QueryRow(ctx, `
						select 
    						appoint_id,
    						doctor_id,
    						patient_id,
							appoint_date,
							appoint_time,
							status,
							to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
							to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
						from appointment where id = $1 `, in.Id).
		Scan(
			&appointment.AppointId,
			&appointment.DoctorId,
			&appointment.PatientId,
			&appointment.AppointDate,
			&appointment.AppointTime,
			&appointment.Status,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

func (a AppointmentRepo) Cancel(ctx context.Context, in *pb.AppointmentId) error {
	_, err := a.Db.Exec(ctx, `update appointment set status = 0 where id = $1`, in.Id)
	if err != nil {
		return err
	}

	return nil
}
