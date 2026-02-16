package middleware

import (
	"electronic-shop-api/config"
	"electronic-shop-api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims représente les données stockées dans le JWT
type Claims struct {
	UserID uint        `json:"user_id"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`
	ShopID uint        `json:"shop_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware vérifie le token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Récupérer le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token d'authentification manquant",
			})
			c.Abort()
			return
		}

		// 2. Vérifier le format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Format du token invalide. Utilisez: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Parser et valider le token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalide ou expiré",
			})
			c.Abort()
			return
		}

		// 4. Stocker les informations dans le contexte
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("shopID", claims.ShopID)

		c.Next()
	}
}

// RequireRole vérifie que l'utilisateur a le rôle requis
func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non authentifié"})
			c.Abort()
			return
		}

		userRole := roleInterface.(models.Role)

		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Accès refusé. Rôle insuffisant.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext récupère les informations utilisateur
func GetUserFromContext(c *gin.Context) (uint, uint, models.Role) {
	userID, _ := c.Get("userID")
	shopID, _ := c.Get("shopID")
	role, _ := c.Get("role")

	return userID.(uint), shopID.(uint), role.(models.Role)
}