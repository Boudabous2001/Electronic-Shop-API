package handlers

import (
	"electronic-shop-api/database"
	"electronic-shop-api/middleware"
	"electronic-shop-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ========================================
// STRUCTURES DE REQUÊTE
// ========================================

type CreateProductInput struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	PurchasePrice float64 `json:"purchase_price" binding:"required,gt=0"`
	SellingPrice  float64 `json:"selling_price" binding:"required,gt=0"`
	Stock         int     `json:"stock" binding:"gte=0"`
	ImageURL      string  `json:"image_url"`
}

type UpdateProductInput struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	Stock         int     `json:"stock"`
	ImageURL      string  `json:"image_url"`
}

// ========================================
// GET ALL PRODUCTS (Private)
// ========================================

func GetProducts(c *gin.Context) {
	_, shopID, role := middleware.GetUserFromContext(c)
	db := database.GetDB()

	var products []models.Product

	// MULTI-TENANT: Filtrer par ShopID du token
	if err := db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}

	// Si Admin, masquer PurchasePrice
	if role == models.RoleAdmin {
		var filteredProducts []gin.H
		for _, p := range products {
			filteredProducts = append(filteredProducts, gin.H{
				"id":            p.ID,
				"name":          p.Name,
				"description":   p.Description,
				"category":      p.Category,
				"selling_price": p.SellingPrice,
				"stock":         p.Stock,
				"image_url":     p.ImageURL,
				"shop_id":       p.ShopID,
				"created_at":    p.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, gin.H{"products": filteredProducts, "count": len(filteredProducts)})
		return
	}

	// SuperAdmin voit tout
	c.JSON(http.StatusOK, gin.H{"products": products, "count": len(products)})
}

// ========================================
// GET SINGLE PRODUCT (Private)
// ========================================

func GetProduct(c *gin.Context) {
	_, shopID, role := middleware.GetUserFromContext(c)

	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	db := database.GetDB()
	var product models.Product

	// MULTI-TENANT: Vérifier que le produit appartient au shop
	if err := db.Where("id = ? AND shop_id = ?", productID, shopID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	// Si Admin, masquer PurchasePrice
	if role == models.RoleAdmin {
		c.JSON(http.StatusOK, gin.H{
			"product": gin.H{
				"id":            product.ID,
				"name":          product.Name,
				"description":   product.Description,
				"category":      product.Category,
				"selling_price": product.SellingPrice,
				"stock":         product.Stock,
				"image_url":     product.ImageURL,
				"shop_id":       product.ShopID,
				"created_at":    product.CreatedAt,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// ========================================
// CREATE PRODUCT
// ========================================

func CreateProduct(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	// Validation: prix de vente > prix d'achat
	if input.SellingPrice < input.PurchasePrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le prix de vente doit être supérieur au prix d'achat"})
		return
	}

	product := models.Product{
		Name:          input.Name,
		Description:   input.Description,
		Category:      input.Category,
		PurchasePrice: input.PurchasePrice,
		SellingPrice:  input.SellingPrice,
		Stock:         input.Stock,
		ImageURL:      input.ImageURL,
		ShopID:        shopID, // Toujours prendre le ShopID du token !
	}

	db := database.GetDB()
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du produit"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produit créé avec succès",
		"product": product,
	})
}

// ========================================
// UPDATE PRODUCT
// ========================================

func UpdateProduct(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()
	var product models.Product

	// MULTI-TENANT
	if err := db.Where("id = ? AND shop_id = ?", productID, shopID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	// Mettre à jour les champs fournis
	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Category != "" {
		updates["category"] = input.Category
	}
	if input.PurchasePrice > 0 {
		updates["purchase_price"] = input.PurchasePrice
	}
	if input.SellingPrice > 0 {
		updates["selling_price"] = input.SellingPrice
	}
	if input.Stock >= 0 {
		updates["stock"] = input.Stock
	}
	if input.ImageURL != "" {
		updates["image_url"] = input.ImageURL
	}

	if err := db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	db.First(&product, productID)
	c.JSON(http.StatusOK, gin.H{"message": "Produit mis à jour", "product": product})
}

// ========================================
// DELETE PRODUCT
// ========================================

func DeleteProduct(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	db := database.GetDB()
	var product models.Product

	// MULTI-TENANT
	if err := db.Where("id = ? AND shop_id = ?", productID, shopID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé avec succès"})
}

// ========================================
// PUBLIC: GET SHOP PRODUCTS (Guest)
// ========================================

func GetPublicProducts(c *gin.Context) {
	shopID, err := strconv.ParseUint(c.Param("shopID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de shop invalide"})
		return
	}

	db := database.GetDB()

	// Vérifier que le shop existe et est actif
	var shop models.Shop
	if err := db.Where("id = ? AND active = ?", shopID, true).First(&shop).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop non trouvé ou inactif"})
		return
	}

	// Récupérer les produits
	var products []models.Product
	query := db.Where("shop_id = ?", shopID)

	// Filtre par catégorie (optionnel)
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Filtre produits en stock (optionnel)
	if inStock := c.Query("in_stock"); inStock == "true" {
		query = query.Where("stock > 0")
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}

	// Convertir en version publique (avec lien WhatsApp)
	publicProducts := make([]models.ProductPublic, 0, len(products))
	for _, p := range products {
		publicProducts = append(publicProducts, p.ToPublic(shop.WhatsAppNumber))
	}

	c.JSON(http.StatusOK, gin.H{
		"shop": gin.H{
			"id":   shop.ID,
			"name": shop.Name,
		},
		"products": publicProducts,
		"count":    len(publicProducts),
	})
}

// ========================================
// PUBLIC: GET SINGLE PRODUCT (Guest)
// ========================================

func GetPublicProduct(c *gin.Context) {
	shopID, err := strconv.ParseUint(c.Param("shopID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de shop invalide"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	db := database.GetDB()

	// Vérifier le shop
	var shop models.Shop
	if err := db.Where("id = ? AND active = ?", shopID, true).First(&shop).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop non trouvé ou inactif"})
		return
	}

	// Récupérer le produit
	var product models.Product
	if err := db.Where("id = ? AND shop_id = ?", productID, shopID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product.ToPublic(shop.WhatsAppNumber)})
}