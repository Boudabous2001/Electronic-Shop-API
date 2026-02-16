package handlers

import (
	"electronic-shop-api/database"
	"electronic-shop-api/middleware"
	"electronic-shop-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ========================================
// GET SHOP INFO
// ========================================

func GetShop(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)
	db := database.GetDB()

	var shop models.Shop
	if err := db.First(&shop, shopID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shop": shop})
}

// ========================================
// UPDATE SHOP
// ========================================

type UpdateShopInput struct {
	Name           string `json:"name"`
	WhatsAppNumber string `json:"whatsapp_number"`
	Active         *bool  `json:"active"`
}

func UpdateShop(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	var input UpdateShopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()
	var shop models.Shop

	if err := db.First(&shop, shopID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop non trouvé"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.WhatsAppNumber != "" {
		updates["whatsapp_number"] = input.WhatsAppNumber
	}
	if input.Active != nil {
		updates["active"] = *input.Active
	}

	if err := db.Model(&shop).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	db.First(&shop, shopID)
	c.JSON(http.StatusOK, gin.H{"message": "Shop mis à jour", "shop": shop})
}

// ========================================
// GET ALL USERS
// ========================================

func GetUsers(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)
	db := database.GetDB()

	var users []models.User
	if err := db.Where("shop_id = ?", shopID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}

	var safeUsers []gin.H
	for _, u := range users {
		safeUsers = append(safeUsers, gin.H{
			"id":         u.ID,
			"name":       u.Name,
			"email":      u.Email,
			"role":       u.Role,
			"shop_id":    u.ShopID,
			"created_at": u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": safeUsers, "count": len(safeUsers)})
}

// ========================================
// CREATE USER
// ========================================

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=SuperAdmin Admin"`
}

func CreateUser(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()

	// Vérifier si l'email existe
	var existingUser models.User
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Cet email est déjà utilisé"})
		return
	}

	// Hasher le mot de passe
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Utilisateur créé",
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
// UPDATE USER
// ========================================

type UpdateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func UpdateUser(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'utilisateur invalide"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()
	var user models.User

	// MULTI-TENANT
	if err := db.Where("id = ? AND shop_id = ?", userID, shopID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Email != "" {
		var existing models.User
		if err := db.Where("email = ? AND id != ?", input.Email, userID).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Cet email est déjà utilisé"})
			return
		}
		updates["email"] = input.Email
	}
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		updates["password"] = string(hashedPassword)
	}
	if input.Role != "" {
		updates["role"] = input.Role
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	db.First(&user, userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Utilisateur mis à jour",
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
// DELETE USER
// ========================================

func DeleteUser(c *gin.Context) {
	userIDFromToken, shopID, _ := middleware.GetUserFromContext(c)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'utilisateur invalide"})
		return
	}

	// Empêcher l'auto-suppression
	if uint(userID) == userIDFromToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vous ne pouvez pas vous supprimer vous-même"})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ? AND shop_id = ?", userID, shopID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur supprimé"})
}

// ========================================
// PUBLIC: GET ALL SHOPS
// ========================================

func GetPublicShops(c *gin.Context) {
	db := database.GetDB()

	var shops []models.Shop
	if err := db.Where("active = ?", true).Find(&shops).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des shops"})
		return
	}

	var publicShops []gin.H
	for _, s := range shops {
		publicShops = append(publicShops, gin.H{
			"id":   s.ID,
			"name": s.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"shops": publicShops, "count": len(publicShops)})
}