package conf

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var dbDefault *mongo.Database

// sử dụng singleton pattern để tạo một connection duy nhất đến database
// khi ứng dụng lớn hơn thì không nên sử dụng singleton pattern
// thay vào đó nên sử dụng connection pool
func (a *App) GetDB() *mongo.Database {
	if dbDefault == nil {
		return a.initDB()
	}
	return dbDefault
}

func (a *App) initDB() *mongo.Database {
	// Thiết lập thông tin kết nối
	clientOptions := options.Client().ApplyURI(cfg.MongoURL)

	// Tạo kết nối đến MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Kiểm tra kết nối
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(cfg.DBName)

	return database
}
