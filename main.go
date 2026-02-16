package main

import (
	"electronic-shop-api/config"
	"electronic-shop-api/database"
	"electronic-shop-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 Démarrage de l'application Electronic Shop API...")

	// Charger la configuration
	config.Load()
	log.Println("✅ Configuration chargée")

	// Connexion à la base de données
	database.Connect()
	log.Println("✅ Base de données connectée")

	// Créer le routeur Gin
	router := gin.Default()

	// Middleware CORS
	router.Use(corsMiddleware())

	// Configurer les routes
	routes.SetupRoutes(router)

	// Démarrer le serveur
	port := config.AppConfig.ServerPort
	log.Printf("🌐 Serveur démarré sur http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Erreur au démarrage du serveur:", err)
	}
}

// corsMiddleware gère les headers CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}