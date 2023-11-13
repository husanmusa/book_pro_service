package storage

import (
	"context"
	pb "github.com/husanmusa/book_pro_service/genproto/book_pro_service"
)

type StorageI interface {
	Book() BookProServiceI
	BookCategory() BookCategoryI
	CloseDB()
}

type BookProServiceI interface {
	CreateBook(ctx context.Context, in *pb.Book) error
	UpdateBook(ctx context.Context, in *pb.Book) error
	GetBookList(ctx context.Context, in *pb.GetBookListRequest) (*pb.GetBookListResponse, error)
	GetBook(ctx context.Context, in *pb.ById) (*pb.Book, error)
	DeleteBook(ctx context.Context, in *pb.ById) error
}

type BookCategoryI interface {
	CreateBookCategory(ctx context.Context, in *pb.BookCategory) error
	UpdateBookCategory(ctx context.Context, in *pb.BookCategory) error
	GetBookCategoryList(ctx context.Context, in *pb.GetBookCategoryListReq) (*pb.GetBookCategoryListRes, error)
	GetBookCategory(ctx context.Context, in *pb.ById) (*pb.BookCategory, error)
	DeleteBookCategory(ctx context.Context, in *pb.ById) error
}
