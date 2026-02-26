package ships

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/dbs/mongodis"
	"github.com/khoirulhasin/untirta_api/app/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type shipMongodistory struct {
	db *mongodis.DB
}

func NewShipMongodistory(db *mongodis.DB) *shipMongodistory {
	return &shipMongodistory{
		db: db,
	}
}

var _ ShipMongodistory = &shipMongodistory{}

func (r *shipMongodistory) GetAllBigShips(ctx context.Context) ([]bson.M, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var source string
	if val := ctx.Value("X-Source"); val != nil {
		if s, ok := val.(string); ok {
			source = s
		}
	}

	// Kunci cache Redis
	cacheKey := "ships:all"
	log.Printf(source)

	// Jika source bukan "crontab", coba ambil dari Redis
	if source != "crontab" {
		cachedData, err := r.db.Redis.Get(timeoutCtx, cacheKey).Result()
		if err == nil {
			var results []bson.M
			if err := json.Unmarshal([]byte(cachedData), &results); err == nil {
				return results, nil
			}
			log.Printf("Failed to unmarshal cached data: %v", err.Error())
		} else if err != redis.Nil {
			log.Printf("Redis error: %v", err.Error())
		}
	}

	// Hanya 3 hari terakhir (gunakan UTC agar konsisten)
	cutoff := time.Now().UTC().Add(-72 * time.Hour)
	filter := bson.M{
		"ts": bson.M{"$gte": cutoff},
	}

	// Sort by timestamp descending
	opts := options.Find().SetSort(bson.D{{Key: "ts", Value: -1}})

	// Query MongoDB dengan filter 3 hari terakhir
	cursor, err := r.db.Mongo.Collection("ais_static").Find(timeoutCtx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeoutCtx)

	var results []bson.M
	if err := cursor.All(timeoutCtx, &results); err != nil {
		return nil, err
	}

	// Jika hasil kosong, simpan ke Redis sebentar untuk mencegah query berulang
	if len(results) == 0 {
		emptyData, _ := json.Marshal([]bson.M{})
		if err := r.db.Redis.Set(timeoutCtx, cacheKey, emptyData, 5*time.Minute).Err(); err != nil {
			log.Printf("Failed to save empty data to Redis: %v", err.Error())
		}
		return nil, nil
	}

	// Simpan hasil ke Redis (TTL 60 menit)
	if data, err := json.Marshal(results); err == nil {
		if err := r.db.Redis.Set(timeoutCtx, cacheKey, data, 60*time.Minute).Err(); err != nil {
			log.Printf("Failed to save data to Redis: %v", err.Error())
		} else {
			log.Printf("Saved %d documents to Redis", len(results))
		}
	} else {
		log.Printf("Failed to marshal data for Redis: %v", err.Error())
	}

	return results, nil
}

func (r *shipMongodistory) GetBigShipsByDatetime(ctx context.Context, durationTimeInput models.DurationTimeInput) ([]bson.M, error) {

	// context timeout untuk Mongo & Redis
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Ambil X-Source (opsional): jika "crontab" => bypass GET cache
	var source string
	if v := ctx.Value("X-Source"); v != nil {
		if s, ok := v.(string); ok {
			source = s
		}
	}

	// Konversi epoch (DETIK). Jika input kamu MILLISECOND, ganti ke time.UnixMilli(...)
	start := time.Unix(int64(durationTimeInput.Start), 0).UTC()
	end := time.Unix(int64(durationTimeInput.End), 0).UTC()

	// Kunci cache (ringkas & deterministik). Tambah versi agar mudah invalidasi di masa depan.
	cacheKey := fmt.Sprintf("ais_dynamic:v1:range:%d-%d", start.Unix(), end.Unix())

	// Coba ambil dari Redis jika bukan crontab
	if source != "crontab" {
		if cached, err := r.db.Redis.Get(timeoutCtx, cacheKey).Result(); err == nil {
			var results []bson.M
			if err := json.Unmarshal([]byte(cached), &results); err == nil {
				return results, nil
			}
			log.Printf("redis unmarshal failed: %v", err)
		} else if err != redis.Nil {
			// selain key-not-found, log error redis
			log.Printf("redis get error: %v", err)
		}
	}

	// Filter waktu end-exclusive biar tidak dobel di tepi 'end'
	filter := bson.M{
		"ts": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	// Sorting + hint indeks (disarankan punya index {timestamp:1} atau {mmsi:1,timestamp:1})
	opts := options.Find().
		SetSort(bson.D{{Key: "ts", Value: -1}}).
		SetHint(bson.D{{Key: "ts", Value: 1}})

	// NOTE: pilih salah satu sesuai wrapper kamu:
	// cur, err := r.db.Collection("ais_dynamic").Find(timeoutCtx, filter, opts)
	cur, err := r.db.Mongo.Collection("ais_dynamic").Find(timeoutCtx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(timeoutCtx)

	var results []bson.M
	if err := cur.All(timeoutCtx, &results); err != nil {
		return nil, err
	}

	// Simpan ke Redis (negative cache jika hasil kosong, TTL lebih pendek)
	if data, err := json.Marshal(results); err == nil {
		ttl := 5 * time.Minute
		if len(results) == 0 {
			ttl = 2 * time.Minute
		}
		if err := r.db.Redis.Set(timeoutCtx, cacheKey, data, ttl).Err(); err != nil {
			log.Printf("redis set error: %v", err)
		}
	} else {
		log.Printf("marshal error for redis: %v", err)
	}

	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}
