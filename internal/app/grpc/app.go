package grpc

import (
	"log/slog"

	"google.golang.org/grpc"

	itemgrpc "item/internal/grpc/item"

)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app.
func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	itemgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}