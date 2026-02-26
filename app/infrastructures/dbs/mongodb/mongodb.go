package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect mengembalikan instance client MongoDB
func Connect() *mongo.Database {
	var (
		databaseUser     = os.Getenv("MONGO_USER")
		databaseName     = os.Getenv("MONGO_NAME")
		databaseHost     = os.Getenv("MONGO_HOST")
		databasePort     = os.Getenv("MONGO_PORT")
		databasePassword = os.Getenv("MONGO_PASS")
	)

	// Format URI untuk MongoDB
	// mongodb://admin:Qazwsxedc.12!!@141.11.241.75:27017/ais_data?authSource=ais_data&authMechanism=SCRAM-SHA-256
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-256",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	// Konfigurasi opsi koneksi
	clientOptions := options.Client().ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	// Koneksi ke MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	// Verifikasi koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Gagal ping MongoDB: %v", err)
	}

	fmt.Println("Berhasil terhubung ke MongoDB!")
	return client.Database(os.Getenv("MONGO_NAME"))
}
