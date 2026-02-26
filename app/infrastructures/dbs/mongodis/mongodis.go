package mongodis

import (
	"github.com/khoirulhasin/untirta_api/app/infrastructures/dbs/mongodb"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/dbs/redisdb"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// DB berisi koneksi ke MongoDB dan Redis
type DB struct {
	Mongo *mongo.Database
	Redis *redis.Client
}

// Connect menginisialisasi koneksi ke MongoDB dan Redis
func Connect() *DB {
	return &DB{
		Mongo: mongodb.Connect(),
		Redis: redisdb.Connect(),
	}
}
