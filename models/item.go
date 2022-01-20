package models

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Item struct {
	ID           uuid.UUID `json:"id"`
	Product      string    `json:"product"`
	Manufacturer string    `json:"manufacturer"`
	DeviceType   string    `json:"device_type"`
	Serial       string    `json:"serial"`
	Condition    string    `json:"condition"`
	Year         string    `json:"year"`
}

func GetItems(db *pgxpool.Pool) ([]Item, error) {
	var items []Item
	err := pgxscan.Select(context.Background(), db, &items, "SELECT * FROM items")
	if err != nil {
		return nil, err
	}

	return items, nil
}

func GetItemsBySearch(db *pgxpool.Pool, search string) ([]Item, error) {
	var items []Item
	err := pgxscan.Select(context.Background(), db, &items, `SELECT * FROM items WHERE to_tsvector(items::text) @@ to_tsquery($1)`, search)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("no items found")
	}

	return items, nil
}

func GetItem(db *pgxpool.Pool, id string) (*Item, error) {
	var item Item
	err := pgxscan.Get(context.Background(), db, &item, "SELECT * FROM items WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func CreateItem(db *pgxpool.Pool, i Item) (*Item, error) {
	var item Item
	err := pgxscan.Get(context.Background(), db, &item, "INSERT INTO items (product, manufacturer, device_type, serial, condition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", &i.Product, &i.Manufacturer, &i.DeviceType, &i.Serial, &i.Condition, &i.Year)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func UpdateItem(db *pgxpool.Pool, id string, i Item) (*Item, error) {
	var item Item
	err := pgxscan.Get(context.Background(), db, &item, "UPDATE items SET product = $2, manufacturer = $3, device_type = $4, serial = $5, condition = $6, year = $7 WHERE id = $1 RETURNING *", id, &i.Product, &i.Manufacturer, &i.DeviceType, &i.Serial, &i.Condition, &i.Year)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func DeleteItem(db *pgxpool.Pool, id string) bool {
	res, err := db.Exec(context.Background(), "DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return false
	}

	count := res.RowsAffected()
	return count != 0
}
