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
// STRUCTURE DE REQUÊTE
// ========================================

type CreateTransactionInput struct {
	Type      string  `json:"type" binding:"required,oneof=Sale Expense Withdrawal"`
	ProductID *uint   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

// ========================================
// GET ALL TRANSACTIONS
// ========================================

func GetTransactions(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)
	db := database.GetDB()

	var transactions []models.Transaction

	query := db.Where("shop_id = ?", shopID).Order("created_at DESC")

	// Filtre par type (optionnel)
	if transType := c.Query("type"); transType != "" {
		query = query.Where("type = ?", transType)
	}

	if err := query.Preload("Product").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions, "count": len(transactions)})
}

// ========================================
// GET SINGLE TRANSACTION
// ========================================

func GetTransaction(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de transaction invalide"})
		return
	}

	db := database.GetDB()
	var transaction models.Transaction

	if err := db.Where("id = ? AND shop_id = ?", transactionID, shopID).
		Preload("Product").First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction non trouvée"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// ========================================
// CREATE TRANSACTION
// ========================================

func CreateTransaction(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	var input CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	db := database.GetDB()

	// ========================================
	// CAS 1: VENTE (Sale)
	// ========================================
	if input.Type == "Sale" {
		// Validation
		if input.ProductID == nil || input.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "product_id et quantity sont requis pour une vente",
			})
			return
		}

		// Vérifier que le produit existe et appartient au shop
		var product models.Product
		if err := db.Where("id = ? AND shop_id = ?", *input.ProductID, shopID).First(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
			return
		}

		// Vérifier le stock
		if product.Stock < input.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":              "Stock insuffisant",
				"stock_disponible":   product.Stock,
				"quantite_demandee":  input.Quantity,
			})
			return
		}

		// Calculer le montant total
		totalAmount := product.SellingPrice * float64(input.Quantity)

		// Transaction DB atomique
		tx := db.Begin()

		// 1. Décrémenter le stock
		if err := tx.Model(&product).Update("stock", product.Stock-input.Quantity).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock"})
			return
		}

		// 2. Créer la transaction
		transaction := models.Transaction{
			Type:      models.TypeSale,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
			Amount:    totalAmount,
			ShopID:    shopID,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la transaction"})
			return
		}

		tx.Commit()

		// Charger le produit pour la réponse
		db.Preload("Product").First(&transaction, transaction.ID)

		c.JSON(http.StatusCreated, gin.H{
			"message":     "Vente enregistrée",
			"transaction": transaction,
			"new_stock":   product.Stock - input.Quantity,
		})
		return
	}

	// ========================================
	// CAS 2: DÉPENSE ou RETRAIT
	// ========================================
	transaction := models.Transaction{
		Type:   models.TransactionType(input.Type),
		Amount: input.Amount,
		ShopID: shopID,
	}

	if err := db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction enregistrée",
		"transaction": transaction,
	})
}

// ========================================
// DELETE TRANSACTION
// ========================================

func DeleteTransaction(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)

	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de transaction invalide"})
		return
	}

	db := database.GetDB()
	var transaction models.Transaction

	if err := db.Where("id = ? AND shop_id = ?", transactionID, shopID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction non trouvée"})
		return
	}

	// Si c'est une vente, restaurer le stock
	if transaction.Type == models.TypeSale && transaction.ProductID != nil {
		var product models.Product
		if err := db.First(&product, *transaction.ProductID).Error; err == nil {
			db.Model(&product).Update("stock", product.Stock+transaction.Quantity)
		}
	}

	if err := db.Delete(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction supprimée"})
}