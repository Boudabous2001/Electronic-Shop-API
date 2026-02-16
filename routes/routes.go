package routes

import (
	"electronic-shop-api/handlers"
	"electronic-shop-api/middleware"
	"electronic-shop-api/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configure toutes les routes de l'API
func SetupRoutes(router *gin.Engine) {

	// ========================================
	// ROUTES PUBLIQUES (Pas d'authentification)
	// ========================================

	// Page d'accueil
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "🏪 Bienvenue sur l'API Electronic Shop",
			"version": "1.0.0",
		})
	})

	// Auth
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Routes publiques pour les clients
	public := router.Group("/public")
	{
		public.GET("/shops", handlers.GetPublicShops)
		public.GET("/:shopID/products", handlers.GetPublicProducts)
		public.GET("/:shopID/products/:productID", handlers.GetPublicProduct)
	}

	// ========================================
	// ROUTES PROTÉGÉES (JWT requis)
	// ========================================

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Profil
		protected.GET("/me", handlers.GetMe)

		// Produits (Admin + SuperAdmin)
		products := protected.Group("/products")
		products.Use(middleware.RequireRole(models.RoleAdmin, models.RoleSuperAdmin))
		{
			products.GET("", handlers.GetProducts)
			products.GET("/:id", handlers.GetProduct)
			products.POST("", handlers.CreateProduct)
			products.PUT("/:id", handlers.UpdateProduct)
			products.DELETE("/:id", handlers.DeleteProduct)
		}

		// Transactions (Admin + SuperAdmin)
		transactions := protected.Group("/transactions")
		transactions.Use(middleware.RequireRole(models.RoleAdmin, models.RoleSuperAdmin))
		{
			transactions.GET("", handlers.GetTransactions)
			transactions.GET("/:id", handlers.GetTransaction)
			transactions.POST("", handlers.CreateTransaction)
			transactions.DELETE("/:id", handlers.DeleteTransaction)
		}

		// Reports (SuperAdmin uniquement)
		reports := protected.Group("/reports")
		reports.Use(middleware.RequireRole(models.RoleSuperAdmin))
		{
			reports.GET("/dashboard", handlers.GetDashboard)
			reports.GET("/low-stock", handlers.GetLowStockProducts)
		}

		// Shop Management (SuperAdmin uniquement)
		shop := protected.Group("/shop")
		shop.Use(middleware.RequireRole(models.RoleSuperAdmin))
		{
			shop.GET("", handlers.GetShop)
			shop.PUT("", handlers.UpdateShop)
		}

		// Users Management (SuperAdmin uniquement)
		users := protected.Group("/users")
		users.Use(middleware.RequireRole(models.RoleSuperAdmin))
		{
			users.GET("", handlers.GetUsers)
			users.POST("", handlers.CreateUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}
	}
}