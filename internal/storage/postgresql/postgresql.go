package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"item-service/internal/config"
	"item-service/internal/domain/models"
	"item-service/pkg/client/postgresql"
)

type Storage struct {
	client postgresql.Client
	log    *slog.Logger
}

func New(log *slog.Logger) *Storage {
	cfg := config.MustLoad()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		log.Info("Failed to connect to PostgreSQL: %v", err)
		return &Storage{}
	}
	log.Info("connected to PostgreSQL")

	return &Storage{
		client: client,
		log:    log,
	}
}

func (s *Storage) SaveItem(ctx context.Context, name string, rarity string, description string) (uuid.UUID, error) {
	const op = "Storage.SaveItem"

	q := `
		INSERT INTO items (
			id, 
			name, 
			rarity, 
			description
		)
		VALUES (
			gen_random_uuid(), 
			$1, 
			$2, 
			$3)
		RETURNING id
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var id uuid.UUID

	if err := s.client.QueryRow(ctx, q, name, rarity, description).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s, op: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState(), op))
			return uuid.Nil, newErr
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)

	}

	s.log.Info("Completed to create item")

	return id, nil
}

func (s *Storage) GetItem(ctx context.Context, itemID uuid.UUID) (*models.Item, error) {
	const op = "Storage.GetItem"

	q := `
        SELECT 
			id, 
			name, 
			rarity, 
			description 
		FROM items
		WHERE id = $1
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var item models.Item
	if err := s.client.QueryRow(ctx, q, itemID).Scan(&item.ItemId, &item.Name, &item.Rarity, &item.Description); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s, op: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState(), op))
			return &models.Item{}, newErr
		}

		return &models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("Completed to get user by id")

	return &item, nil
}

func (s *Storage) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	const op = "Storage.GetAllItems"

	q := `
        SELECT 
			id, 
			name, 
			rarity, 
			description 
		FROM items
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s, op: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState(), op))
			return []*models.Item{}, newErr
		}

		return []*models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	items := make([]*models.Item, 0)

	for rows.Next() {
		var item models.Item

		if err := rows.Scan(&item.ItemId, &item.Name, &item.Rarity, &item.Description); err != nil {
			return []*models.Item{}, fmt.Errorf("%s: %w", op, err)
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return []*models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}

func (s *Storage) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	const op = "Storage.DeleteItem"

	q := `
		DELETE FROM items
		WHERE id = $1
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := s.client.Exec(ctx, q, itemID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s, op: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState(), op))

			return newErr
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
