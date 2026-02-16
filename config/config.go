package config

import (
	"os"
	"time"
)

// Config contient toutes les configurations de l'application
type Config struct {
	JWTSecret     string
	JWTExpiration time.Duration
	ServerPort    string
}

// AppConfig est l'instance globale de configuration
var AppConfig Config

// Load charge les configurations
func Load() {
	AppConfig = Config{
		// Clé secrète pour signer les JWT
		JWTSecret: getEnv("JWT_SECRET", "votre-cle-secrete-super-securisee-2024"),

		// Durée de validité du token (24 heures)
		JWTExpiration: 24 * time.Hour,

		// Port du serveur
		ServerPort: getEnv("PORT", "8080"),
	}
}

// getEnv récupère une variable d'environnement ou retourne une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}