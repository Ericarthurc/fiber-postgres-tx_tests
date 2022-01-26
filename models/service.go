package models

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	ID    uuid.UUID `json:"id"`
	Time  time.Time `json:"time"`
	Seats int       `json:"seats"`
}

func GetServices(db *pgxpool.Pool) ([]Service, error) {
	var Services []Service
	err := pgxscan.Select(context.Background(), db, &Services, "SELECT * FROM Services")
	if err != nil {
		return nil, err
	}

	return Services, nil
}

func GetService(db *pgxpool.Pool, id string) (*Service, error) {
	var Service Service
	err := pgxscan.Get(context.Background(), db, &Service, "SELECT * FROM Services WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &Service, nil
}

func CreateService(db *pgxpool.Pool, i Service) (*Service, error) {
	var Service Service
	err := pgxscan.Get(context.Background(), db, &Service, "INSERT INTO Services (product, manufacturer, device_type, serial, condition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", &i.Product, &i.Manufacturer, &i.DeviceType, &i.Serial, &i.Condition, &i.Year)
	if err != nil {
		return nil, err
	}

	return &Service, nil
}

func UpdateService(db *pgxpool.Pool, id string, i Service) (*Service, error) {
	var Service Service
	err := pgxscan.Get(context.Background(), db, &Service, "UPDATE Services SET product = $2, manufacturer = $3, device_type = $4, serial = $5, condition = $6, year = $7 WHERE id = $1 RETURNING *", id, &i.Product, &i.Manufacturer, &i.DeviceType, &i.Serial, &i.Condition, &i.Year)
	if err != nil {
		return nil, err
	}

	return &Service, nil
}

func DeleteService(db *pgxpool.Pool, id string) bool {
	res, err := db.Exec(context.Background(), "DELETE FROM Services WHERE id = $1", id)
	if err != nil {
		return false
	}

	count := res.RowsAffected()
	return count != 0
}
