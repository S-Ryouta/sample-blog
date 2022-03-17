package db

import (
	"github.com/S-Ryouta/sample-blog/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	DBConn *gorm.DB
)

func Connect() *gorm.DB {
	user := os.Getenv("DATABASE_USERNAME")
	pass := os.Getenv("DATABASE_PASSWORD")
	dsn := user + ":" + pass + "@tcp(sample-blog-db:3306)/" + "sample_blog" + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	tableMigrate(db)
	DBConn = db
	return DBConn
}

func tableMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.Entry{})
}
