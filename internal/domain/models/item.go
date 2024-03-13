package models

import "github.com/google/uuid"

type Item struct {
	ItemId  uuid.UUID `json:"item_id"`
	Name    string    `json:"name" validate:"required,min=3,max=100"`
	Rarity  string    `json:"rarity" validate:"required,min=3,max=20"`
	Quality string    `json:"quality,omitempty" validate:"required,min=3,max=1000"`
}
