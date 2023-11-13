package service

import (
	"context"
	pb "github.com/husanmusa/book_pro_service/genproto/book_pro_service"
	"github.com/husanmusa/book_pro_service/pkg/helper"
	"github.com/husanmusa/book_pro_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bookService struct {
	logger  logger.LoggerI
	storage storage.StorageI
	pb.UnimplementedBookProServiceServer
}

func NewBookService(log logger.LoggerI, store storage.StorageI) *bookService {
	return &bookService{
		logger:  log,
		storage: store,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req *pb.Book) (*emptypb.Empty, error) {
	s.logger.Info("---Create Book--->", logger.Any("req", req))
	err := s.storage.Book().CreateBook(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error creating book", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *bookService) GetBook(ctx context.Context, req *pb.ById) (*pb.Book, error) {
	s.logger.Info("---Get Book--->", logger.Any("req", req))
	resp, err := s.storage.Book().GetBook(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting book", req, codes.Internal)
	}
	return resp, nil
}

func (s *bookService) GetBookList(ctx context.Context, req *pb.GetBookListRequest) (*pb.GetBookListResponse, error) {
	s.logger.Info("---Get Book--->", logger.Any("req", req))
	resp, err := s.storage.Book().GetBookList(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting book", req, codes.Internal)
	}
	return resp, nil
}

func (s *bookService) UpdateBook(ctx context.Context, req *pb.Book) (*emptypb.Empty, error) {
	s.logger.Info("---Update Book--->", logger.Any("req", req))
	err := s.storage.Book().UpdateBook(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error updating book", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *bookService) DeleteBook(ctx context.Context, req *pb.ById) (*emptypb.Empty, error) {
	s.logger.Info("---Delete Book--->", logger.Any("req", req))
	err := s.storage.Book().DeleteBook(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error deleting book", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}
