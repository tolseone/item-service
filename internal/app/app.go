package app

import (
	"log/slog"

	grpcapp "item-service/internal/app/grpc"
	"item-service/internal/service"
	db "item-service/internal/storage/postgresql"

)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int) *App {
	storage := db.New(log)
	if storage == nil {
		panic("Failed to create storage")
	}

	itemService := item.New(log, storage)

	grpcApp := grpcapp.New(log, itemService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
