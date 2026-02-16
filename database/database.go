package database

import (
	"electronic-shop-api/models"
	"log"

	"github.com/glebarez/sqlite" // Driver SQLite pure Go (pas de CGO)
	"gorm.io/gorm"
)

// DB est l'instance globale de la base de données
var DB *gorm.DB

// Connect initialise la connexion à la base de données
func Connect() {
	var err error

	// Connexion SQLite avec le driver pure Go
	DB, err = gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Échec de connexion à la base de données:", err)
	}

	log.Println("✅ Connecté à la base de données SQLite")

	// Migration automatique des tables
	err = DB.AutoMigrate(
		&models.Shop{},
		&models.User{},
		&models.Product{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("❌ Échec de migration:", err)
	}

	log.Println("✅ Migration des tables terminée")
}

// GetDB retourne l'instance de la base de données
func GetDB() *gorm.DB {
	return DB
}