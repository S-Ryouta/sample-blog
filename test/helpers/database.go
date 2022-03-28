package helpers

import (
	"github.com/S-Ryouta/sample-blog/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
)

type DbConfig struct {
	user         string
	pass         string
	host         string
	databaseName string
}

type TestDb interface {
	SetUp()
	Connect(useDatabaseName bool)
	CleanUp()
}

func InitDb() *DbConfig {
	var config *DbConfig
	config = &DbConfig{
		user:         os.Getenv("TEST_DATABASE_USERNAME"),
		pass:         os.Getenv("TEST_DATABASE_PASSWORD"),
		host:         os.Getenv("TEST_DATABASE_HOST"),
		databaseName: os.Getenv("TEST_DATABASE_NAME"),
	}

	return config
}

func (config *DbConfig) Connect(useDatabaseName bool) *gorm.DB {
	var dsn string
	if useDatabaseName {
		dsn = config.user + ":" + config.pass + "@tcp(" + config.host + ":3306)/" + config.databaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dsn = config.user + ":" + config.pass + "@tcp(" + config.host + ":3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected " + dsn)

	return db
}

func (config *DbConfig) Setup() {
	db := config.Connect(false)

	db = db.Exec("CREATE DATABASE IF NOT EXISTS " + config.databaseName)
	mysqlDB, _ := db.DB()
	defer mysqlDB.Close()
}

func (config *DbConfig) CleanUp() {
	db := config.Connect(false)

	db = db.Exec("DROP DATABASE " + config.databaseName)
	mysqlDB, _ := db.DB()
	defer mysqlDB.Close()
}

func TableMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.Entry{})
}

func TestDbMock(t *testing.T) {
	initDb := InitDb()
	t.Setenv("DATABASE_NAME", initDb.databaseName)
	t.Setenv("DATABASE_HOST", initDb.host)
}
