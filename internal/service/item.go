package item

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"item-service/internal/domain/models"
	"item-service/internal/lib/logger/sl"
	"item-service/internal/storage"
)

type Item struct {
	log  *slog.Logger
	repo RepositoryItem
}

type RepositoryItem interface {
	SaveItem(ctx context.Context, name string, rarity string, description string) (itemID uuid.UUID, err error)
	DeleteItem(ctx context.Context, itemID uuid.UUID) (err error)
	GetAllItems(ctx context.Context) (items []*models.Item, err error)
	GetItem(ctx context.Context, itemID uuid.UUID) (item *models.Item, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// New returns a new instance of the Item service.
func New(log *slog.Logger, repo RepositoryItem) *Item {
	return &Item{
		repo: repo,
		log:  log,
	}
}

// CreateItem creates a new item.
func (itm *Item) CreateItem(ctx context.Context, name, rarity, description string) (uuid.UUID, error) {
	const op = "Item.CreateItem"

	log := itm.log.With(
		slog.String("op", op),
		slog.String("name", name),
		slog.String("rarity", rarity),
		slog.String("description", description),
	)

	log.Info("attempting to create item")

	itemID, err := itm.repo.SaveItem(ctx, name, rarity, description)
	if err != nil {
		if errors.Is(err, storage.ErrItemExists) {
			itm.log.Warn("item already exists", sl.Err(err))

			return uuid.Nil, fmt.Errorf("$s: %w", ErrInvalidCredentials)
		}

		itm.log.Error("failed to create item", sl.Err(err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	itm.log.Info("item successfully created")

	return itemID, nil
}

// GetItem returns the item with the given ID.
func (itm *Item) GetItem(ctx context.Context, itemID uuid.UUID) (*models.Item, error) {
	const op = "Item.GetItem"

	log := itm.log.With(
		slog.String("op", op),
		slog.Any("itemID", itemID),
	)

	log.Info("attemting to get item")

	item, err := itm.repo.GetItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			log.Warn("item not found", sl.Err(err))

			return &models.Item{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get item", sl.Err(err))

		return &models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("item successfully got")

	return item, nil
}

// GetAllItems returns all items.
func (itm *Item) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	const op = "Item.GetAllItems"

	log := itm.log.With(
		slog.String("op", op),
	)

	log.Info("attemting to get all items")

	items, err := itm.repo.GetAllItems(ctx)
	if err != nil {
		itm.log.Info("failed to get all items", sl.Err(err))

		return []*models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("All items received")

	return items, nil
}

// DeleteItem deletes the item with the given ID.
func (itm *Item) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	const op = "Item.GetAllItems"

	log := itm.log.With(
		slog.String("op", op),
	)

	log.Info("attemting to get all items")

	if err := itm.repo.DeleteItem(ctx, itemID); err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			itm.log.Warn("items not found", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
		itm.log.Info("failed to get all items", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Item successfully deleted")

	return nil
}
