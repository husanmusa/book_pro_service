package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/husanmusa/book_pro_service/genproto/book_pro_service"
	"github.com/husanmusa/book_pro_service/pkg/helper"
	"github.com/husanmusa/book_pro_service/storage"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BookCategoryRepo struct {
	Db *pgxpool.Pool
}

func NewBookCategoryRepo(db *pgxpool.Pool) storage.BookCategoryI {
	return &BookCategoryRepo{
		Db: db,
	}
}

func (r *BookCategoryRepo) CreateBookCategory(ctx context.Context, req *pb.BookCategory) error {
	_, err := r.Db.Exec(ctx, `insert into book_category (
					id,
                	name
                	) 
		values ($1, $2)`,
		uuid.NewString(),
		req.Name,
	)
	if err != nil {
		return fmt.Errorf("error while create book, err: %s", err.Error())
	}

	return nil
}

func (r *BookCategoryRepo) GetBookCategoryList(ctx context.Context, req *pb.GetBookCategoryListReq) (*pb.GetBookCategoryListRes, error) {
	var (
		bookcts []*pb.BookCategory
		count   int32
		arr     []interface{}
		params  = make(map[string]interface{})
	)

	offset := ""
	filter := " where true "
	limit := "  "

	if req.Offset > 0 {
		offset = " OFFSET :offset "
		params["offset"] = req.Offset
	}

	if req.Limit > 0 {
		limit = " LIMIT :limit "
		params["limit"] = req.Limit
	}

	cQ := `select count(1) from book_category `
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	err := r.Db.QueryRow(ctx, cQ, arr...).Scan(
		&count,
	)
	if err != nil {
		return nil, err
	}

	q := `select 
			id,
			name
			from book_category ` + filter + offset + limit

	q, arr = helper.ReplaceQueryParams(q, params)
	rows, err := r.Db.Query(ctx, q, arr...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			bookct pb.BookCategory
		)

		err = rows.Scan(
			&bookct.Id,
			&bookct.Name,
		)
		if err != nil {
			return nil, err
		}

		bookcts = append(bookcts, &bookct)
	}

	return &pb.GetBookCategoryListRes{
		BookCategories: bookcts,
		Count:          count,
	}, nil
}

func (r *BookCategoryRepo) GetBookCategory(ctx context.Context, req *pb.ById) (*pb.BookCategory, error) {
	var (
		book pb.BookCategory
	)

	err := r.Db.QueryRow(ctx, `select 
			id,
			name
from book_category where id = $1 `,
		req.Id).Scan(
		&book.Id,
		&book.Name,
	)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookCategoryRepo) UpdateBookCategory(ctx context.Context, req *pb.BookCategory) error {
	query := `update book set
                name=$1,
              where id=$2 `

	_, err := r.Db.Exec(ctx, query, req.Name, req.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *BookCategoryRepo) DeleteBookCategory(ctx context.Context, req *pb.ById) error {
	query := `
		delete from book_category
		where id = $1`

	result, err := r.Db.Exec(ctx, query, req.Id)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
