package postgres

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/husanmusa/med-appointment-service/genproto/appointment_service"
	"github.com/husanmusa/med-appointment-service/pkg/helper"
	"github.com/husanmusa/med-appointment-service/storage"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PatientRepo struct {
	Db *pgxpool.Pool
}

func NewPatientRepo(db *pgxpool.Pool) storage.PatientI {
	return &PatientRepo{
		Db: db,
	}
}

func (r *PatientRepo) Create(ctx context.Context, req *pb.Patient) error {
	_, err := r.Db.Exec(ctx, `insert into patient (
					id,
                	name,
                	age,
                	gender,
                	login,
                	password,
                	access_token,
                	refresh_token
                    ) 
		values ($1, $2, $3, $4, $5, $6, '', '')`,
		uuid.NewString(),
		req.Name,
		req.Age,
		req.Gender,
		req.Login,
		req.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PatientRepo) GetList(ctx context.Context, req *pb.GetListPatientsRequest) (*pb.GetListPatientsResponse, error) {
	var (
		patients []*pb.Patient
		count    int32
		arr      []interface{}
		params   = make(map[string]interface{})
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

	cQ := `select count(1) from patient `
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
			gender,
			login
		from patient
		where true ` + order + offset + limit

	q, arr = helper.ReplaceQueryParams(q, params)
	rows, err := r.Db.Query(ctx, q, arr...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			patient pb.Patient
		)

		err = rows.Scan(
			&patient.PatientId,
			&patient.Name,
			&patient.Age,
			&patient.Gender,
			&patient.Login,
		)
		if err != nil {
			return nil, err
		}

		patients = append(patients, &patient)
	}

	return &pb.GetListPatientsResponse{
		Patients: patients,
		Count:    count,
	}, nil
}

func (r *PatientRepo) Get(ctx context.Context, req *pb.PatientId) (*pb.Patient, error) {
	var (
		patient pb.Patient
	)

	err := r.Db.QueryRow(ctx, `select 
    				id,
					name,
					age,
					gender,
					login
			from patient where id = $1 `,
		req.PatientId).Scan(
		&patient.PatientId,
		&patient.Name,
		&patient.Age,
		&patient.Gender,
		&patient.Login,
	)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *PatientRepo) Update(ctx context.Context, req *pb.Patient) error {

	query := `update patient set
                name=$1,
			  	age=$2,
				gender=$3,
				login=$4,
				updated_at=current_timestamp
              where id=$5 `

	_, err := r.Db.Exec(ctx, query, req.Name, req.Age, req.Gender, req.Login, req.PatientId)

	if err != nil {
		return err
	}

	return nil
}

func (r *PatientRepo) Delete(ctx context.Context, req *pb.PatientId) error {
	filter := " WHERE "
	params := make(map[string]interface{})

	params["id"] = req.PatientId
	filter += " id = :id "

	query := "delete from patient " + filter
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
