package item

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	itemv1 "github.com/tolseone/protos/gen/go/item"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"item-service/internal/domain/models"

)

type Item interface {
	CreateItem(ctx context.Context, name, rarity, description string) (itemID uuid.UUID, err error)
	GetItem(ctx context.Context, itemID uuid.UUID) (item *models.Item, err error)
	GetAllItems(ctx context.Context) (items []*models.Item, err error)
	DeleteItem(ctx context.Context, itemID uuid.UUID) (err error)
}

type serverAPI struct {
	itemv1.UnimplementedItemServiceServer
	item      Item
	validator *validator.Validate
}

func Register(gRPC *grpc.Server, item Item) {
	validator := validator.New()
	itemv1.RegisterItemServiceServer(gRPC, &serverAPI{
		item:      item,
		validator: validator,
	})
}

func (s *serverAPI) CreateItem(ctx context.Context, req *itemv1.CreateItemRequest) (*itemv1.CreateItemResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	itemID, err := s.item.CreateItem(ctx, req.Item.GetName(), req.Item.GetRarity(), req.Item.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	itemIDString := itemID.String()
	if itemIDString == "" {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &itemv1.CreateItemResponse{
		ItemId: itemIDString,
	}, nil
}

func (s *serverAPI) GetItem(ctx context.Context, req *itemv1.GetItemRequest) (*itemv1.GetItemResponse, error) {
	panic("implement me")
}

func (s *serverAPI) GetAllItems(ctx context.Context, req *itemv1.GetAllItemsRequest) (*itemv1.GetAllItemsResponse, error) {
	panic("implement me")
}

func (s *serverAPI) DeleteItem(ctx context.Context, req *itemv1.DeleteItemRequest) (*itemv1.DeleteItemResponse, error) {
	panic("implement me")
}
