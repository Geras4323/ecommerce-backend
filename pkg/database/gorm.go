package database

import (
	"log"

	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
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

	dsn := utils.GetEnvVar("DB_DSN")
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
		&models.CartItem{},
		&models.Product{},
		&models.Image{},
		&models.Order{},
		&models.OrderProduct{},
		&models.Payment{},
		&models.State{},
		&models.Unit{},
	)

	if migrationError != nil {
		log.Fatal("DATABASE: error during migration: ", migrationError.Error())
	}
}
