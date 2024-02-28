package models

import "github.com/google/uuid"

type Item struct {
	ItemId      uuid.UUID `json:"item_id"`
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	Rarity      string    `json:"rarity" validate:"required,min=3,max=20"`
	Description string    `json:"description,omitempty" validate:"required,min=3,max=1000"`
}

type CreateItemRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Rarity      string `json:"rarity" validate:"required,min=3,max=20"`
	Description string `json:"description,omitempty" validate:"required,min=3,max=1000"`
}

type GetItemRequest struct {
	ItemId uuid.UUID `json:"item_id" validate:"required"`
}

type DeleteItemRequest struct {
	ItemId uuid.UUID `json:"item_id" validate:"required"`
}
