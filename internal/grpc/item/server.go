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
	CreateItem(ctx context.Context, name, rarity, quality string) (itemID uuid.UUID, err error)
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

	itemID, err := s.item.CreateItem(ctx, req.GetName(), req.GetRarity(), req.GetQuality())
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
	if err := s.validator.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	itemID, err := uuid.Parse(req.GetItemId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse item id")
	}

	item, err := s.item.GetItem(ctx, itemID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get item")
	}

	return &itemv1.GetItemResponse{
		Item: &itemv1.Item{
			ItemId:  item.ItemId.String(),
			Name:    item.Name,
			Rarity:  item.Rarity,
			Quality: item.Quality,
		},
	}, nil
}

func (s *serverAPI) GetAllItems(ctx context.Context, req *itemv1.GetAllItemsRequest) (*itemv1.GetAllItemsResponse, error) {
	items, err := s.item.GetAllItems(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all items")
	}

	var itemResponses []*itemv1.Item
	for _, item := range items {
		itemResponses = append(itemResponses, &itemv1.Item{
			ItemId:  item.ItemId.String(),
			Name:    item.Name,
			Rarity:  item.Rarity,
			Quality: item.Quality,
		})
	}

	response := &itemv1.GetAllItemsResponse{
		Items: itemResponses,
	}

	if len(itemResponses) == 0 {
		return &itemv1.GetAllItemsResponse{}, nil
	}

	return response, nil
}

func (s *serverAPI) DeleteItem(ctx context.Context, req *itemv1.DeleteItemRequest) (*itemv1.DeleteItemResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	itemID, err := uuid.Parse(req.GetItemId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse item id")
	}

	if err := s.item.DeleteItem(ctx, itemID); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete item")
	}

	return &itemv1.DeleteItemResponse{}, nil
}
