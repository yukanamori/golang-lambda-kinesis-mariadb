package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DBConfig struct {
		User     string
		Password string
		Host     string
		Port     string
		DBName   string
	}

	Record struct {
		gorm.Model
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Field3 bool   `json:"field3"`
	}
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	db, err := initDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}

	for _, record := range kinesisEvent.Records {
		data := record.Kinesis.Data
		log.Println("Received data:", string(data))

		var r Record
		if err := json.Unmarshal(data, &r); err != nil {
			log.Println("Error unmarshalling data:", err)
			return err
		}

		result := db.Create(&r)
		if err := result.Error; err != nil {
			log.Println("Error inserting data into database:", err)
			return err
		}
	}

	return nil
}

func initDB() (*gorm.DB, error) {
	dbConfig := DBConfig{
		User:     "user",
		Password: "password",
		Host:     "host",
		Port:     "port",
		DBName:   "dbname",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Record{}); err != nil {
		return nil, err
	}

	return db, nil
}
