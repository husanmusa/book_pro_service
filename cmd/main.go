package main

import (
	"context"
	"fmt"
	"github.com/husanmusa/book_pro_service/config"
	pb "github.com/husanmusa/book_pro_service/genproto/book_pro_service"
	"github.com/husanmusa/book_pro_service/service"
	"github.com/husanmusa/book_pro_service/storage/postgres"
	"log"
	"net"

	"github.com/saidamir98/udevs_pkg/logger"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var loggerLevel string
	cfg := config.Load()

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
	case config.TestMode:
		loggerLevel = logger.LevelDebug
	default:
		loggerLevel = logger.LevelInfo
	}

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Panic("error while listen RPC", logger.Error(err))
		return
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer func() {
		if err := logger.Cleanup(log); err != nil {
			log.Error("Failed to cleanup logger", logger.Error(err))
		}
	}()

	maxMsgSize := 100 * 1024 * 1024

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	pgStore, err := postgres.NewPostgres(context.Background(), cfg, log)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer pgStore.CloseDB()

	bookService := service.NewBookService(log, pgStore)
	bookCategoryService := service.NewBookCategoryService(log, pgStore)
	reflection.Register(s)

	pb.RegisterBookProServiceServer(s, bookService)
	pb.RegisterBookCategoryServiceServer(s, bookCategoryService)

	group, ctx := errgroup.WithContext(context.Background())
	fmt.Println(ctx)

	group.Go(func() error {
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
		log.Panic("Api server has finished")
		return nil
	})

	err = group.Wait()
	if err != nil {
		log.Panic("error while listening: %v", logger.Error(err))
	}
}
