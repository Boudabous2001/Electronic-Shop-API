package handlers

import (
	"electronic-shop-api/config"
	"electronic-shop-api/database"
	"electronic-shop-api/middleware"
	"electronic-shop-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ========================================
// STRUCTURES DE REQUÊTE
// ========================================

type RegisterInput struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	Role           string `json:"role" binding:"required,oneof=SuperAdmin Admin"`
	ShopName       string `json:"shop_name"`
	ShopID         uint   `json:"shop_id"`
	WhatsAppNumber string `json:"whatsapp_number"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ========================================
// REGISTER
// ========================================

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()

	// Vérifier si l'email existe déjà
	var existingUser models.User
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Cet email est déjà utilisé"})
		return
	}

	// Gérer le Shop
	var shopID uint

	if input.ShopID != 0 {
		// Shop existant
		var shop models.Shop
		if err := db.First(&shop, input.ShopID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shop non trouvé"})
			return
		}
		shopID = shop.ID
	} else {
		// Nouveau shop
		if input.ShopName == "" || input.WhatsAppNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "shop_name et whatsapp_number sont requis pour créer un nouveau shop",
			})
			return
		}

		newShop := models.Shop{
			Name:           input.ShopName,
			Active:         true,
			WhatsAppNumber: input.WhatsAppNumber,
		}

		if err := db.Create(&newShop).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du shop"})
			return
		}
		shopID = newShop.ID
	}

	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hashage du mot de passe"})
		return
	}

	// Créer l'utilisateur
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     models.Role(input.Role),
		ShopID:   shopID,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	// Générer le token
	token, _ := generateToken(user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Compte créé avec succès",
		"token":   token,
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    user.Role,
			"shop_id": user.ShopID,
		},
	})
}

// ========================================
// LOGIN
// ========================================

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()

	// Trouver l'utilisateur
	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Vérifier le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Vérifier que le shop est actif
	var shop models.Shop
	if err := db.First(&shop, user.ShopID).Error; err != nil || !shop.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ce shop est désactivé"})
		return
	}

	// Générer le token
	token, _ := generateToken(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Connexion réussie",
		"token":   token,
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    user.Role,
			"shop_id": user.ShopID,
		},
	})
}

// ========================================
// GENERATE TOKEN
// ========================================

func generateToken(user models.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		ShopID: user.ShopID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ========================================
// GET ME
// ========================================

func GetMe(c *gin.Context) {
	userID, shopID, role := middleware.GetUserFromContext(c)

	db := database.GetDB()
	var user models.User
	db.First(&user, userID)

	var shop models.Shop
	db.First(&shop, shopID)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    role,
			"shop_id": shopID,
		},
		"shop": gin.H{
			"id":              shop.ID,
			"name":            shop.Name,
			"whatsapp_number": shop.WhatsAppNumber,
		},
	})
}