package database

import (
	"log"

	"github.com/geras4323/ecommerce-backend/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Gorm *gorm.DB

func ConnectGorm() {
	if Gorm != nil {
		log.Println("DATABASE: already connected")
		return
	}

	var dialector gorm.Dialector

	dsn := "geras:admin123@tcp(127.0.0.1:3305)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	dialector = mysql.Open(dsn)

	conn, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatal("DATABASE: failed to connect")
	}

	Gorm = conn

	migrationError := conn.AutoMigrate(
		&models.Category{},
		&models.Supplier{},
		&models.User{},
		&models.Product{},
		&models.Image{},
		&models.Order{},
		&models.OrderProduct{},
		&models.Payment{},
	)

	if migrationError != nil {
		log.Fatal("DATABASE: error during migration: ", migrationError.Error())
	}
}
