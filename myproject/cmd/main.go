package main

import (
	"context"
	"log"
	"myproject/internal/bookstore"
	"myproject/internal/config"
	"myproject/internal/handlers"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := bookstore.NewPostgresDatabase(cfg.GetConnectionString())
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	if db != nil {
		defer db.Close()
	}

	bs := bookstore.NewBookStore(db)
	h := handlers.NewBookHandlers(bs)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := db.Ping(); err != nil {
				log.Printf("Database connection lost: %v", err)
				// พยายามเชื่อมต่อใหม่
				if reconnErr := db.Reconnect(cfg.GetConnectionString()); reconnErr != nil {
					log.Printf("Failed to reconnect: %v", reconnErr)
				} else {
					log.Printf("Successfully reconnected to the database")
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(TimeoutMiddleware(5 * time.Second))

	r.GET("/health", h.HealthCheck)

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/AllStoreInfo", h.GetAllStoreInfo)
		v1.GET("/store/:id", h.GetStoreInfoByID)
		v1.GET("/product/:store_id", h.GetProductsByStore)
		v1.GET("/newproduct/:store_id", h.GetNewProductsByStore)
		v1.GET("/searchproducts", h.SearchProducts)
		v1.GET("/products/:id", h.GetProduct)
		v1.GET("/:store_id/search", h.SearchProductsByStore)
		v1.GET("/Allproduct/:store_id/sort", h.GetAllProductsByStore)
		v1.GET("/:store_id/by-category", h.GetProductsByCategoryAndStore)
		v1.GET("/category", h.GetALLProductsByCategory)
		v1.POST("/store/:store_id/product/:product_id/add_to_cart", h.AddToCart)
		v1.GET("/cart/:store_id", h.GetCartItemsByStore) // เพิ่มเส้นทางนี้
		v1.DELETE("/store/:store_id/product/:product_id/remove_from_cart", h.DeleteProductFromCart)
		v1.GET("/all-guitars", h.GetAllGuitars)

		// เส้นทางสำหรับ Checkout (ย้ายข้อมูลจาก cart ไป order_history)
		v1.POST("/checkout/:store_id", h.Checkout)
	}

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Printf("Failed to run server: %v", err)
	}
}
