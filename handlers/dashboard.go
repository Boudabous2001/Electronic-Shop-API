package handlers

import (
	"electronic-shop-api/database"
	"electronic-shop-api/middleware"
	"electronic-shop-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ========================================
// GET DASHBOARD (SuperAdmin uniquement)
// ========================================

func GetDashboard(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)
	db := database.GetDB()

	// 1. Total des ventes
	var totalSales float64
	db.Model(&models.Transaction{}).
		Where("shop_id = ? AND type = ?", shopID, models.TypeSale).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSales)

	// 2. Total des dépenses
	var totalExpenses float64
	db.Model(&models.Transaction{}).
		Where("shop_id = ? AND type = ?", shopID, models.TypeExpense).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpenses)

	// 3. Total des retraits
	var totalWithdrawals float64
	db.Model(&models.Transaction{}).
		Where("shop_id = ? AND type = ?", shopID, models.TypeWithdrawal).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalWithdrawals)

	// 4. Coût des produits vendus
	var costOfGoodsSold float64
	db.Raw(`
		SELECT COALESCE(SUM(p.purchase_price * t.quantity), 0)
		FROM transactions t
		JOIN products p ON t.product_id = p.id
		WHERE t.shop_id = ? AND t.type = ?
	`, shopID, models.TypeSale).Scan(&costOfGoodsSold)

	// 5. Profit net
	netProfit := totalSales - costOfGoodsSold - totalExpenses

	// 6. Produits en stock faible (< 5)
	var lowStockCount int64
	db.Model(&models.Product{}).
		Where("shop_id = ? AND stock < ?", shopID, 5).
		Count(&lowStockCount)

	// 7. Total des produits
	var totalProducts int64
	db.Model(&models.Product{}).
		Where("shop_id = ?", shopID).
		Count(&totalProducts)

	// 8. Valeur totale du stock
	var stockValue float64
	db.Model(&models.Product{}).
		Where("shop_id = ?", shopID).
		Select("COALESCE(SUM(purchase_price * stock), 0)").
		Scan(&stockValue)

	// 9. Nombre de transactions par type
	var salesCount, expensesCount, withdrawalsCount int64
	db.Model(&models.Transaction{}).Where("shop_id = ? AND type = ?", shopID, models.TypeSale).Count(&salesCount)
	db.Model(&models.Transaction{}).Where("shop_id = ? AND type = ?", shopID, models.TypeExpense).Count(&expensesCount)
	db.Model(&models.Transaction{}).Where("shop_id = ? AND type = ?", shopID, models.TypeWithdrawal).Count(&withdrawalsCount)

	// 10. Top 5 produits vendus
	type TopProduct struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		TotalSold   int     `json:"total_sold"`
		TotalAmount float64 `json:"total_amount"`
	}
	var topProducts []TopProduct
	db.Raw(`
		SELECT 
			t.product_id,
			p.name as product_name,
			SUM(t.quantity) as total_sold,
			SUM(t.amount) as total_amount
		FROM transactions t
		JOIN products p ON t.product_id = p.id
		WHERE t.shop_id = ? AND t.type = ?
		GROUP BY t.product_id, p.name
		ORDER BY total_sold DESC
		LIMIT 5
	`, shopID, models.TypeSale).Scan(&topProducts)

	c.JSON(http.StatusOK, gin.H{
		"dashboard": gin.H{
			"total_sales":        totalSales,
			"total_expenses":     totalExpenses,
			"total_withdrawals":  totalWithdrawals,
			"cost_of_goods_sold": costOfGoodsSold,
			"net_profit":         netProfit,
			"gross_margin":       totalSales - costOfGoodsSold,
			"total_products":     totalProducts,
			"low_stock_products": lowStockCount,
			"stock_value":        stockValue,
			"transactions": gin.H{
				"sales":       salesCount,
				"expenses":    expensesCount,
				"withdrawals": withdrawalsCount,
				"total":       salesCount + expensesCount + withdrawalsCount,
			},
			"top_products": topProducts,
		},
	})
}

// ========================================
// GET LOW STOCK PRODUCTS
// ========================================

func GetLowStockProducts(c *gin.Context) {
	_, shopID, _ := middleware.GetUserFromContext(c)
	db := database.GetDB()

	var products []models.Product
	if err := db.Where("shop_id = ? AND stock < ?", shopID, 5).
		Order("stock ASC").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
		"message":  "Produits avec stock inférieur à 5 unités",
	})
}