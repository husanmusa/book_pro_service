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

type BookServiceRepo struct {
	Db *pgxpool.Pool
}

func NewBookServiceRepo(db *pgxpool.Pool) storage.BookProServiceI {
	return &BookServiceRepo{
		Db: db,
	}
}

func (r *BookServiceRepo) CreateBook(ctx context.Context, req *pb.Book) error {
	_, err := r.Db.Exec(ctx, `insert into book (
					id,
                	name,
                	author_name,
                	pages,
					description,
                	book_category_id) 
		values ($1, $2, $3, $4, $5, $6)`,
		uuid.NewString(),
		req.Name,
		req.AuthorName,
		req.Pages,
		req.Description,
		req.BookCategoryId,
	)
	if err != nil {
		return fmt.Errorf("error while create book, err: %s", err.Error())
	}

	return nil
}

func (r *BookServiceRepo) GetBookList(ctx context.Context, req *pb.GetBookListRequest) (*pb.GetBookListResponse, error) {
	var (
		books  []*pb.Book
		count  int32
		arr    []interface{}
		params = make(map[string]interface{})
	)

	offset := ""
	filter := " where true "
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

	if len(req.BookCategoryId) > 0 {
		filter += " and b.book_category_id = :book_category_id "
		params["book_category_id"] = req.BookCategoryId
	}

	cQ := `select count(1) from book as b left join book_category as bc on b.book_category_id = bc.id` + filter
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	err := r.Db.QueryRow(ctx, cQ, arr...).Scan(
		&count,
	)
	if err != nil {
		return nil, err
	}

	q := `select 
			b.id,
			b.name,
			b.author_name,
			b.pages,
			b.description,
			bc.name,
			to_char(b.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
from book as b left join book_category as bc on b.book_category_id = bc.id ` + filter + order + offset + limit
	q, arr = helper.ReplaceQueryParams(q, params)

	rows, err := r.Db.Query(ctx, q, arr...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			book pb.Book
		)

		err = rows.Scan(
			&book.Id,
			&book.Name,
			&book.AuthorName,
			&book.Pages,
			&book.Description,
			&book.BookCategory,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	return &pb.GetBookListResponse{
		Books: books,
		Count: count,
	}, nil
}

func (r *BookServiceRepo) GetBook(ctx context.Context, req *pb.ById) (*pb.Book, error) {
	var (
		book pb.Book
	)

	err := r.Db.QueryRow(ctx, `select 
			b.id,
			b.name,
			b.author_name,
			b.pages,
			b.description,
			bc.name,
			to_char(b.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
from book as b left join book_category as bc on b.book_category_id = bc.id where b.id = $1 `,
		req.Id).Scan(
		&book.Id,
		&book.Name,
		&book.AuthorName,
		&book.Pages,
		&book.Description,
		&book.BookCategory,
		&book.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookServiceRepo) UpdateBook(ctx context.Context, req *pb.Book) error {
	query := `update book set
                name=$1,
			  	author_name=$2,
				pages=$3,
				description=$4,
				book_category_id=$5,
				updated_at=current_timestamp
              where id=$6 `

	_, err := r.Db.Exec(ctx, query, req.Name, req.AuthorName, req.Pages, req.Description, req.BookCategoryId, req.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *BookServiceRepo) DeleteBook(ctx context.Context, req *pb.ById) error {
	query := `
		UPDATE book SET
			deleted_at = date_part('epoch', CURRENT_TIMESTAMP)::int
		WHERE id = $1 and deleted_at = 0
	`

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
