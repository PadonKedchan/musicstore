// book_handlers.go
package handlers

import (
	"myproject/internal/bookstore"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetAllGuitars ดึงสินค้าทุกชิ้นที่มีคำว่า "กีตาร์"
func (h *BookHandlers) GetAllGuitars(c *gin.Context) {
	var allGuitars []bookstore.Product

	// ดึงข้อมูลร้านทั้งหมด
	stores, err := h.bs.GetAllStoreInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ลูปผ่านทุกร้านเพื่อดึงสินค้าจากร้านนั้นๆ
	for _, store := range stores {
		// ดึงสินค้าทั้งหมดจากร้านนี้
		products, err := h.bs.GetProductsByStore(c.Request.Context(), store.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// กรองสินค้าที่ชื่อมีคำว่า "กีตาร์"
		for _, product := range products {
			if strings.Contains(product.ProductName, "กีตาร์") {
				allGuitars = append(allGuitars, product)
			}
		}
	}

	// ถ้าไม่มีสินค้ากีตาร์ในร้าน
	if len(allGuitars) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No guitars found"})
		return
	}

	// ส่งข้อมูลสินค้ากีตาร์ทั้งหมดกลับไป
	c.JSON(http.StatusOK, gin.H{"guitars": allGuitars})
}

type BookHandlers struct {
	bs *bookstore.BookStore
}

func NewBookHandlers(bs *bookstore.BookStore) *BookHandlers {
	return &BookHandlers{bs: bs}
}

func (h *BookHandlers) HealthCheck(c *gin.Context) {
	err := h.bs.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"reason": "Database connection failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (h *BookHandlers) GetAllStoreInfo(c *gin.Context) {
	stores, err := h.bs.GetAllStoreInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"store_info": stores})
}

func (h *BookHandlers) GetStoreInfoByID(c *gin.Context) {
	// ดึง ID จาก URL parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// เรียกฟังก์ชัน GetStoreInfoByID จาก BookStore
	store, err := h.bs.GetStoreInfoByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// ส่งข้อมูลร้านในรูปแบบ JSON
	c.JSON(http.StatusOK, gin.H{
		"id":           store.ID,
		"logo_path":    store.LogoPath,
		"store_name":   store.StoreName,
		"description":  store.Description,
		"address":      store.Address,
		"phone_number": store.PhoneNumber,
		"email":        store.Email,
	})
}

func (h *BookHandlers) GetProductsByStore(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	products, err := h.bs.GetProductsByStore(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found for this store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"store_id": storeID, "products": products})
}

func (h *BookHandlers) GetNewProductsByStore(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	products, err := h.bs.GetNewProductsByStore(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found for this store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"store_id": storeID, "products": products})
}

func (h *BookHandlers) SearchProducts(c *gin.Context) {
	// รับค่าพารามิเตอร์จาก URL query string
	productName := c.DefaultQuery("product_name", "") // ใช้ DefaultQuery เพื่อตั้งค่าดีฟอลต์เมื่อไม่มีข้อมูล

	// ถ้าไม่มีค่าของ product_name
	if productName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// ค้นหาผลิตภัณฑ์ที่มีชื่อตรงกับ productName
	products, err := h.bs.SearchProducts(c.Request.Context(), productName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่งข้อมูลผลิตภัณฑ์ที่ค้นหากลับไป
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *BookHandlers) GetProduct(c *gin.Context) {
	// รับค่าจาก URL parameter :id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// ค้นหาผลิตภัณฑ์
	product, err := h.bs.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// ส่งข้อมูลสินค้า
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *BookHandlers) SearchProductsByStore(c *gin.Context) {
	// ดึง store_id จาก URL parameter
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// รับค่าพารามิเตอร์จาก URL query string
	productName := c.DefaultQuery("product_name", "")

	// ตรวจสอบว่ามีการส่งคำค้นหาหรือไม่
	if productName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// เรียกใช้ฟังก์ชันค้นหาผลิตภัณฑ์ในร้านที่กำหนดจาก BookStore
	products, err := h.bs.SearchProductsByStore(c.Request.Context(), productName, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่งข้อมูลผลิตภัณฑ์ที่ค้นหากลับไป
	c.JSON(http.StatusOK, gin.H{"store_id": storeID, "products": products})
}

func (h *BookHandlers) GetAllProductsByStore(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// รับค่า sortOrder จาก query string ถ้ามี
	sortOrder := c.DefaultQuery("sortOrder", "asc") // ค่าเริ่มต้นเป็น "asc"

	// เรียกฟังก์ชัน GetAllProductsByStore จาก BookStore พร้อมกับ sortOrder
	products, err := h.bs.GetAllProductsByStore(c.Request.Context(), storeID, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found for this store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"store_id": storeID, "products": products})
}

// ตัวอย่างสำหรับ Go (Gin framework)
func (h *BookHandlers) GetProductsByCategoryAndStore(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	}

	products, err := h.bs.GetProductsByCategoryAndStore(c.Request.Context(), storeID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found for this store and category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"store_id": storeID, "category": category, "products": products})
}

func (h *BookHandlers) GetALLProductsByCategory(c *gin.Context) {
	category := c.DefaultQuery("category", "") // รับค่าหมวดหมู่จาก query string

	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	}

	// เรียกฟังก์ชัน GetALLProductsByCategory จาก BookStore
	products, err := h.bs.GetALLProductsByCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found for this category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category, "products": products})
}

func (h *BookHandlers) AddToCart(c *gin.Context) {
	// ตรวจสอบ store_id
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// ตรวจสอบ product_id
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// ตรวจสอบ quantity และตั้งค่าเริ่มต้นเป็น 1 ถ้าไม่ได้ส่งมา
	quantityStr := c.DefaultPostForm("quantity", "1")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// เพิ่มสินค้าลงในตะกร้า
	err = h.bs.AddToCart(c.Request.Context(), storeID, productID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่ง response ว่าสินค้าได้ถูกเพิ่มในตะกร้า
	c.JSON(http.StatusOK, gin.H{
		"message":    "Product added to cart",
		"store_id":   storeID,
		"product_id": productID,
		"quantity":   quantity,
	})
}

func (h *BookHandlers) GetCartItemsByStore(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// ดึงข้อมูลสินค้าที่อยู่ในตะกร้าของร้านนั้น
	products, err := h.bs.GetCartItemsByStore(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products in the cart for this store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"store_id": storeID, "cart_items": products})
}

func (h *BookHandlers) DeleteProductFromCart(c *gin.Context) {
	// ตรวจสอบ store_id
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// ตรวจสอบ product_id
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// เรียกใช้ฟังก์ชันลบสินค้าจากตะกร้าใน BookStore
	err = h.bs.DeleteProductFromCart(c.Request.Context(), storeID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่ง response ว่าลบสินค้าจากตะกร้าเรียบร้อยแล้ว
	c.JSON(http.StatusOK, gin.H{
		"message":    "Product removed from cart",
		"store_id":   storeID,
		"product_id": productID,
	})
}

func (h *BookHandlers) Checkout(c *gin.Context) {
	// ตรวจสอบ store_id
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// เรียกใช้ฟังก์ชันดึงสินค้าทั้งหมดในตะกร้าของร้านนั้น
	cartItems, err := h.bs.GetCartItemsByStore(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items in cart to checkout"})
		return
	}

	// คำนวณยอดรวมการสั่งซื้อจากสินค้าที่อยู่ในตะกร้า
	totalAmount := 0.0
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// ทดสอบการชำระเงิน (ในที่นี้เป็นแค่การจำลองการทำงาน)
	paymentSuccess := true // ควรเปลี่ยนให้เป็นการตรวจสอบจากระบบชำระเงินจริง ๆ

	if !paymentSuccess {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "Payment failed"})
		return
	}

	// อัปเดตสถานะสินค้าในตะกร้าให้เป็น 'checked_out'
	err = h.bs.CheckoutCart(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่ง response ว่าการสั่งซื้อเสร็จสมบูรณ์
	c.JSON(http.StatusOK, gin.H{
		"message":      "Checkout successful",
		"store_id":     storeID,
		"total_amount": totalAmount,
	})
}
