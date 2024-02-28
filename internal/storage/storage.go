package storage

import "errors"

var (
	ErrItemExists   = errors.New("Item already exists")
	ErrItemNotFound = errors.New("Item not found")
)
