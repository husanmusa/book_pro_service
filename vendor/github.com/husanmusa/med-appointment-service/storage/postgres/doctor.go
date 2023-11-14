package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/husanmusa/med-appointment-service/genproto/appointment_service"
	"github.com/husanmusa/med-appointment-service/pkg/helper"
	"github.com/husanmusa/med-appointment-service/storage"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DoctorRepo struct {
	Db *pgxpool.Pool
}

func NewDoctorRepo(db *pgxpool.Pool) storage.DoctorI {
	return &DoctorRepo{
		Db: db,
	}
}

func (r *DoctorRepo) Create(ctx context.Context, req *pb.Doctor) error {
	_, err := r.Db.Exec(ctx, `insert into doctor (
					id,
                	name,
                	age,
                	role,
                	polyclinic,
                	gender) 
		values ($1, $2, $3, $4, $5, $6)`,
		uuid.NewString(),
		req.Name,
		req.Age,
		req.Role,
		req.Polyclinic,
		req.Gender,
	)
	if err != nil {
		return fmt.Errorf("error while create doctor, err: %s", err.Error())
	}

	return nil
}

func (r *DoctorRepo) GetList(ctx context.Context, req *pb.GetListDoctorsRequest) (*pb.GetListDoctorsResponse, error) {
	var (
		doctors []*pb.Doctor
		count   int32
		arr     []interface{}
		params  = make(map[string]interface{})
	)

	offset := ""
	limit := " LIMIT 10 "
	order := " ORDER BY created_at DESC "

	if req.Offset > 0 {
		offset = " OFFSET :offset "
		params["offset"] = req.Offset
	}

	if req.Limit > 0 {
		limit = " LIMIT :limit "
		params["limit"] = req.Limit
	}

	cQ := `select count(1) from doctor `
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	err := r.Db.QueryRow(ctx, cQ, arr...).Scan(
		&count,
	)
	if err != nil {
		return nil, err
	}

	q := `select 
			id,
			name,
			age,
			role,
			polyclinic,
			gender
		from doctor
		where true ` + order + offset + limit

	q, arr = helper.ReplaceQueryParams(q, params)
	rows, err := r.Db.Query(ctx, q, arr...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			doctor pb.Doctor
		)

		err = rows.Scan(
			&doctor.DoctorId,
			&doctor.Name,
			&doctor.Age,
			&doctor.Role,
			&doctor.Polyclinic,
			&doctor.Gender,
		)
		if err != nil {
			return nil, err
		}

		doctors = append(doctors, &doctor)
	}

	return &pb.GetListDoctorsResponse{
		Doctors: doctors,
		Count:   count,
	}, nil
}

func (r *DoctorRepo) Get(ctx context.Context, req *pb.DoctorId) (*pb.Doctor, error) {
	var (
		doctor pb.Doctor
	)

	err := r.Db.QueryRow(ctx, `select 
    				id,
					name,
					age,
					role,
					polyclinic,
					gender
			from doctor where id = $1 `,
		req.DoctorId).Scan(
		&doctor.DoctorId,
		&doctor.Name,
		&doctor.Age,
		&doctor.Role,
		&doctor.Polyclinic,
		&doctor.Gender,
	)
	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepo) Update(ctx context.Context, req *pb.Doctor) error {

	query := `update doctor set
                name=$1,
			  	age=$2,
				role=$3,
				polyclinic=$4,
				gender=$5,
				updated_at=current_timestamp
              where id=$6 `

	_, err := r.Db.Exec(ctx, query, req.Name, req.Age, req.Role, req.Polyclinic, req.Gender, req.DoctorId)

	if err != nil {
		return err
	}

	return nil
}

func (r *DoctorRepo) Delete(ctx context.Context, req *pb.DoctorId) error {
	filter := " WHERE "
	params := make(map[string]interface{})

	params["id"] = req.DoctorId
	filter += " id = :id "

	query := "delete from doctor " + filter
	query, arr := helper.ReplaceQueryParams(query, params)

	result, err := r.Db.Exec(ctx, query, arr...)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
