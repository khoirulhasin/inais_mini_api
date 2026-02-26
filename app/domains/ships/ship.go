package ships

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ShipRepository interface {
	CreateShip(ctx context.Context, country *models.Ship) (*models.Ship, error)
	UpdateShip(ctx context.Context, id int32, country *models.Ship) (*models.Ship, error)
	UpdateShipByUUID(ctx context.Context, uuid string, country *models.Ship) (*models.Ship, error)
	DeleteShip(ctx context.Context, id int32) error
	DeleteShipByUUID(ctx context.Context, uuid string) error
	GetShipByID(ctx context.Context, id int32) (*models.Ship, error)
	GetShipByUUID(ctx context.Context, uuid string) (*models.Ship, error)
	GetAllShips(ctx context.Context) ([]*models.Ship, error)
	PageShip(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}

type ShipMongotory interface {
	GetShipsByImei(imei string, durationTimeInput models.DurationTimeInput) ([]bson.M, error)
	GetShipsByDatetime(durationTimeInput models.DurationTimeInput, mmsiList []int64) ([]bson.M, error)
	GetMobShips(durationTimeInput models.DurationTimeInput) ([]bson.M, error)
}

type ShipMongodistory interface {
	GetAllBigShips(ctx context.Context) ([]bson.M, error)
	GetBigShipsByDatetime(ctx context.Context, durationTimeInput models.DurationTimeInput) ([]bson.M, error)
}
