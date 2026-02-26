package ships

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/khoirulhasin/untirta_api/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type shipMongotory struct {
	db *mongo.Database
}

func NewShipMongotory(db *mongo.Database) *shipMongotory {
	return &shipMongotory{
		db,
	}
}

var _ ShipMongotory = &shipMongotory{}

func (r *shipMongotory) GetShipsByImei(imei string, durationTimeInput models.DurationTimeInput) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Misalnya unix time dalam detik
	startUnix := int64(durationTimeInput.Start) // contoh
	endUnix := int64(durationTimeInput.End)     // contoh

	// Konversi ke time.Time (UTC)
	start := time.Unix(startUnix, 0).UTC()
	end := time.Unix(endUnix, 0).UTC()

	// Format jadi "2025-09-14T21:42:07.984+00:00"
	// RFC3339Nano akan menghasilkan format dengan Z, ganti jadi +00:00
	startStr := start.Format("2006-01-02T15:04:05.000+00:00")
	endStr := end.Format("2006-01-02T15:04:05.000+00:00")
	log.Print(startStr)
	// Define the filter to find documents by IMEI and timestamp range
	filter := bson.M{
		"imei": imei,
		"ts": bson.M{
			"$gte": startStr,
			"$lte": endStr,
		},
	}

	// Define sort option by timestamp ascending
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}) // 1 = ascending, -1 = descending

	// Query the "records" collection with sorting
	cursor, err := r.db.Collection("ais_dynamic").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func (r *shipMongotory) GetShipsByDatetime(durationTimeInput models.DurationTimeInput, mmsiList []int64) ([]bson.M, error) {
	// Validasi dasar
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Unix(int64(durationTimeInput.Start), 0).UTC()
	end := time.Unix(int64(durationTimeInput.End), 0).UTC()

	// Define the filter to find documents by MMSI array, timestamp range
	// Also ensure latitude and longitude exist and are not null
	filter := bson.M{
		"ts": bson.M{
			"$gte": start,
			"$lte": end,
		},
		"decoded.Latitude": bson.M{
			"$exists": true,
			"$ne":     nil,
		},
		"decoded.Longitude": bson.M{
			"$exists": true,
			"$ne":     nil,
		},
	}
	log.Print(mmsiList)
	// Tambahkan filter MMSI jika array tidak kosong
	if len(mmsiList) > 0 {
		filter["mmsi"] = bson.M{"$in": mmsiList}
	}

	// Define sort option by timestamp descending
	opts := options.Find().SetSort(bson.D{{"ts", -1}}) // 1 = ascending, -1 = descending

	// Query the "ais_dynamic" collection with sorting
	cursor, err := r.db.Collection("ais_dynamic").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

// func (r *shipMongotory) GetMobShips(durationTimeInput models.DurationTimeInput) ([]bson.M, error) {
// 	// Validasi dasar
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
// 	defer cancel()

// 	start := time.Unix(int64(durationTimeInput.Start), 0).UTC()
// 	end := time.Unix(int64(durationTimeInput.End), 0).UTC()
// 	log.Print(start)

// 	// Menggunakan aggregation pipeline untuk mendapatkan data unique berdasarkan mmsi dan text
// 	// dengan mengambil data terakhir (berdasarkan timestamp descending)
// 	pipeline := []bson.M{
// 		// Stage 1: Filter by timestamp range
// 		{
// 			"$match": bson.M{
// 				"ts": bson.M{
// 					"$gte": start,
// 					"$lte": end,
// 				},
// 			},
// 		},
// 		// Stage 2: Sort by timestamp descending (terbaru dulu)
// 		{
// 			"$sort": bson.M{
// 				"mmsi":         1,
// 				"decoded.Text": 1,
// 				"ts":           -1,
// 			},
// 		},
// 		// Stage 3: Group by mmsi and text, ambil yang pertama (terakhir karena sudah di-sort desc)
// 		{
// 			"$group": bson.M{
// 				"_id": bson.M{
// 					"mmsi": "$mmsi",
// 					"text": "$decoded.Text",
// 				},
// 				"doc": bson.M{"$first": "$$ROOT"}, // Mengambil dokumen pertama dari setiap group
// 			},
// 		},
// 		// Stage 4: Replace root untuk mengembalikan struktur dokumen asli
// 		{
// 			"$replaceRoot": bson.M{
// 				"newRoot": "$doc",
// 			},
// 		},
// 		// Stage 5: Sort final result by timestamp descending
// 		{
// 			"$sort": bson.M{
// 				"ts": -1,
// 			},
// 		},
// 	}

// 	// Execute aggregation
// 	opts := options.Aggregate().SetAllowDiskUse(true)
// 	cursor, err := r.db.Collection("ais_mob").Aggregate(ctx, pipeline, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var results []bson.M
// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, err
// 	}

// 	if len(results) == 0 {
// 		return nil, nil
// 	}

// 	return results, nil
// }

func (r *shipMongotory) GetMobShips(durationTimeInput models.DurationTimeInput) ([]bson.M, error) {
	// Validasi dasar
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Unix(int64(durationTimeInput.Start), 0).UTC()
	end := time.Unix(int64(durationTimeInput.End), 0).UTC()
	log.Print(start)

	// Pipeline untuk message type 14 (MOB) - group by mmsi dan text
	mobPipeline := []bson.M{
		{
			"$match": bson.M{
				"ts":           bson.M{"$gte": start, "$lte": end},
				"message_type": 14,
			},
		},
		{
			"$sort": bson.M{
				"mmsi":         1,
				"decoded.Text": 1,
				"ts":           -1,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"mmsi": "$mmsi",
					"text": "$decoded.Text",
				},
				"doc": bson.M{"$first": "$$ROOT"},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$doc",
			},
		},
	}

	// Pipeline untuk message type 1 (Position) - group by mmsi saja
	positionPipeline := []bson.M{
		{
			"$match": bson.M{
				"ts":           bson.M{"$gte": start, "$lte": end},
				"message_type": 1,
			},
		},
		{
			"$sort": bson.M{
				"mmsi": 1,
				"ts":   -1,
			},
		},
		{
			"$group": bson.M{
				"_id": "$mmsi",
				"doc": bson.M{"$first": "$$ROOT"},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$doc",
			},
		},
	}

	// Execute aggregations
	opts := options.Aggregate().SetAllowDiskUse(true)

	// Get MOB data
	mobCursor, err := r.db.Collection("ais_mob").Aggregate(ctx, mobPipeline, opts)
	if err != nil {
		return nil, err
	}
	defer mobCursor.Close(ctx)

	var mobResults []bson.M
	if err := mobCursor.All(ctx, &mobResults); err != nil {
		return nil, err
	}

	// Get Position data
	posCursor, err := r.db.Collection("ais_mob").Aggregate(ctx, positionPipeline, opts)
	if err != nil {
		return nil, err
	}
	defer posCursor.Close(ctx)

	var posResults []bson.M
	if err := posCursor.All(ctx, &posResults); err != nil {
		return nil, err
	}

	// Gabungkan semua hasil
	var allResults []bson.M
	allResults = append(allResults, mobResults...)
	allResults = append(allResults, posResults...)

	// Sort berdasarkan timestamp terbaru
	sort.Slice(allResults, func(i, j int) bool {
		tsI, okI := allResults[i]["ts"].(primitive.DateTime)
		tsJ, okJ := allResults[j]["ts"].(primitive.DateTime)

		if !okI || !okJ {
			return false
		}

		return tsI.Time().After(tsJ.Time())
	})

	if len(allResults) == 0 {
		return nil, nil
	}

	return allResults, nil
}
