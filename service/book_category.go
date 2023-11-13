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

type bookCategoryService struct {
	logger  logger.LoggerI
	storage storage.StorageI
	pb.UnimplementedBookCategoryServiceServer
}

func NewBookCategoryService(log logger.LoggerI, store storage.StorageI) *bookCategoryService {
	return &bookCategoryService{
		logger:  log,
		storage: store,
	}
}

func (s *bookCategoryService) CreateBookCategory(ctx context.Context, req *pb.BookCategory) (*emptypb.Empty, error) {
	s.logger.Info("---Create BookCategory--->", logger.Any("req", req))
	err := s.storage.BookCategory().CreateBookCategory(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error creating bookCategory", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *bookCategoryService) GetBookCategory(ctx context.Context, req *pb.ById) (*pb.BookCategory, error) {
	s.logger.Info("---Get BookCategory--->", logger.Any("req", req))
	resp, err := s.storage.BookCategory().GetBookCategory(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting bookCategory", req, codes.Internal)
	}
	return resp, nil
}

func (s *bookCategoryService) GetBookCategoryList(ctx context.Context, req *pb.GetBookCategoryListReq) (*pb.GetBookCategoryListRes, error) {
	s.logger.Info("---Get BookCategory--->", logger.Any("req", req))
	resp, err := s.storage.BookCategory().GetBookCategoryList(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error getting bookCategory", req, codes.Internal)
	}
	return resp, nil
}

func (s *bookCategoryService) UpdateBookCategory(ctx context.Context, req *pb.BookCategory) (*emptypb.Empty, error) {
	s.logger.Info("---Update BookCategory--->", logger.Any("req", req))
	err := s.storage.BookCategory().UpdateBookCategory(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error updating bookCategory", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

func (s *bookCategoryService) DeleteBookCategory(ctx context.Context, req *pb.ById) (*emptypb.Empty, error) {
	s.logger.Info("---Delete BookCategory--->", logger.Any("req", req))
	err := s.storage.BookCategory().DeleteBookCategory(ctx, req)
	if err != nil {
		return nil, helper.HandleError(s.logger, err, "Error deleting bookCategory", req, codes.Internal)
	}
	return &emptypb.Empty{}, nil
}
