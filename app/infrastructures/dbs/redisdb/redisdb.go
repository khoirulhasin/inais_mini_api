package redisdb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect returns a Redis Client instance
func Connect() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST") // Format: host:port, e.g., "localhost:6379"
	redisPassword := os.Getenv("REDIS_PASSWORD")

	fmt.Printf("REDIS_HOST: %s\n", redisHost)
	fmt.Printf("REDIS_PASSWORD: %s\n", redisPassword)

	if redisHost == "" {
		log.Fatal("REDIS_HOST environment variable is not set")
	}

	// Configure Redis options for a standalone instance
	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		Password:     redisPassword, // Empty string if no password is set
		DB:           0,             // Default DB
		DialTimeout:  30 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis!")
	return client
}
