// bookstore.go
package bookstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// โครงสร้างสำหรับเก็บข้อมูลหนังสือจากทุกฟิลด์ที่ต้องการ
type Book struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Author      *string `json:"author"` // ใช้ *string เพื่อรองรับค่า NULL
	Image       *string `json:"image"`
	Description *string `json:"description"`
	DateAdded   *string `json:"created_at"`
	Category    *string `json:"category"`
}

type StoreInfo struct {
	ID          int    `json:"id"`
	LogoPath    string `json:"logo_path"`
	StoreName   string `json:"store_name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type Product struct {
	ID            int       `json:"id"`
	ProductName   string    `json:"product_name"`
	Price         float64   `json:"price"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Category      string    `json:"category"`
	Brand         string    `json:"brand"`
	Model         string    `json:"model"`
	StoreID       int       `json:"store_id"`
	IsRecommended bool      `json:"is_recommended"`
	ImagePath     string    `json:"image_path"`
}

// BookDatabase เป็น Interface ที่กำหนดว่า Book Database ต้องทำอะไรได้บ้าง
type BookDatabase interface {
	GetAllStoreInfo(ctx context.Context) ([]StoreInfo, error) // เพิ่มฟังก์ชันนี้
	Close() error
	Ping() error
	Reconnect(connStr string) error
	GetStoreInfoByID(ctx context.Context, id int) (StoreInfo, error)
	GetProductsByStore(ctx context.Context, storeID int) ([]Product, error)
	GetNewProductsByStore(ctx context.Context, storeID int) ([]Product, error)
	SearchProducts(ctx context.Context, searchQuery string) ([]Product, error)
	GetProduct(ctx context.Context, id int) (Product, error)
	SearchProductsByStore(ctx context.Context, searchQuery string, storeID int) ([]Product, error)
	GetAllProductsByStore(ctx context.Context, storeID int, sortOrder string) ([]Product, error)
	GetProductsByCategoryAndStore(ctx context.Context, storeID int, category string) ([]Product, error)
	GetALLProductsByCategory(ctx context.Context, category string) ([]Product, error)
	AddToCart(ctx context.Context, storeID, productID, quantity int) error
	GetCartItemsByStore(ctx context.Context, storeID int) ([]Product, error)
	DeleteProductFromCart(ctx context.Context, storeID, productID int) error
	Checkout(ctx context.Context, storeID int) error
	CheckoutCart(ctx context.Context, storeID int) error
}

// PostgresDatabase เป็น struct ที่เชื่อมต่อกับ PostgreSQL Database จริง
type PostgresDatabase struct {
	db *sql.DB
}

// NewPostgresDatabase สร้าง PostgresDatabase ใหม่และเชื่อมต่อกับฐานข้อมูล
func NewPostgresDatabase(connStr string) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// ตั้งค่า connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ทดสอบการเชื่อมต่อด้วย context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresDatabase{db: db}, nil
}

func (pdb *PostgresDatabase) Close() error {
	return pdb.db.Close()
}

func (pdb *PostgresDatabase) Ping() error {
	if pdb == nil {
		return errors.New("database connection is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pdb.db.PingContext(ctx)
}

func (pdb *PostgresDatabase) Reconnect(connStr string) error {
	if pdb.db != nil {
		pdb.db.Close()
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// ตั้งค่า connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	pdb.db = db
	return nil
}

// BookStore เป็นโครงสร้างหลักของ Application
type BookStore struct {
	db BookDatabase
}

// NewBookStore สร้าง BookStore ใหม่โดยรับ Database ที่จะใช้
func NewBookStore(db BookDatabase) *BookStore {
	return &BookStore{db: db}
}

// Close เป็น Method ของ BookStore ที่ใช้ปิดการเชื่อมต่อกับฐานข้อมูล
func (bs *BookStore) Close() error {
	return bs.db.Close()
}

func (bs *BookStore) Ping() error {
	if bs.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return bs.db.Ping()
}

func (pdb *PostgresDatabase) GetAllStoreInfo(ctx context.Context) ([]StoreInfo, error) {
	query := `SELECT id, logo_path, store_name, description, address, phone_number, email FROM store_info`
	rows, err := pdb.db.QueryContext(ctx, query) // ใช้ pdb.db ซึ่งเป็น *sql.DB
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []StoreInfo
	for rows.Next() {
		var store StoreInfo
		if err := rows.Scan(&store.ID,
			&store.LogoPath,
			&store.StoreName,
			&store.Description,
			&store.Address,
			&store.PhoneNumber,
			&store.Email); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}

func (pdb *PostgresDatabase) GetStoreInfoByID(ctx context.Context, id int) (StoreInfo, error) {
	var store StoreInfo
	query := `SELECT id, logo_path, store_name, description, address, phone_number, email FROM store_info WHERE id = $1`
	err := pdb.db.QueryRowContext(ctx, query, id).Scan(
		&store.ID,
		&store.LogoPath,
		&store.StoreName,
		&store.Description,
		&store.Address,
		&store.PhoneNumber,
		&store.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return store, fmt.Errorf("store not found")
		}
		return store, fmt.Errorf("failed to get store: %v", err)
	}
	return store, nil
}

// เพิ่มฟังก์ชันใน PostgresDatabase สำหรับการดึงข้อมูลสินค้าจาก store_id
func (pdb *PostgresDatabase) GetProductsByStore(ctx context.Context, storeID int) ([]Product, error) {
	query := `SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path
                FROM product_info WHERE store_id = $1 ORDER BY created_at DESC LIMIT 3;;`
	rows, err := pdb.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (pdb *PostgresDatabase) GetNewProductsByStore(ctx context.Context, storeID int) ([]Product, error) {
	query := `SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path 
				FROM product_info WHERE store_id = $1 ORDER BY created_at desc LIMIT 1;`
	rows, err := pdb.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (pdb *PostgresDatabase) SearchProducts(ctx context.Context, searchQuery string) ([]Product, error) {
	// Query ที่จะค้นหาผลิตภัณฑ์ที่ตรงกับคำค้นหาในบางส่วน
	query := `
        SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model,  store_id, is_recommended, image_path 
        FROM product_info 
        WHERE product_name ILIKE $1 
        ORDER BY product_name
    `

	// ใช้ '%' เพื่อให้ค้นหาคำที่มีตัวอักษรตรงส่วนใดส่วนหนึ่ง เช่น 'P' จะเจอ 'phone'
	rows, err := pdb.db.QueryContext(ctx, query, "%"+searchQuery+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

// แสดงสินค้า 1 อัน
func (pdb *PostgresDatabase) GetProduct(ctx context.Context, id int) (Product, error) {
	var product Product
	// แก้ไข query เพื่อให้ตรงกับตารางและฟิลด์ของ Product
	err := pdb.db.QueryRowContext(ctx, `
        SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path 
        FROM product_info WHERE id = $1`, id).Scan(
		&product.ID,
		&product.ProductName,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Category,
		&product.Brand,
		&product.Model,
		&product.StoreID,
		&product.IsRecommended,
		&product.ImagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, fmt.Errorf("product not found")
		}
		return product, fmt.Errorf("failed to get product: %v", err)
	}
	return product, nil
}

func (pdb *PostgresDatabase) SearchProductsByStore(ctx context.Context, searchQuery string, storeID int) ([]Product, error) {
	// Query ที่จะค้นหาผลิตภัณฑ์ที่ตรงกับคำค้นหาในบางส่วนและเฉพาะร้านที่กำหนด
	query := `
        SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path 
        FROM product_info 
        WHERE product_name ILIKE $1 
        AND store_id = $2
        ORDER BY product_name
    `

	// ใช้ '%' เพื่อให้ค้นหาคำที่มีตัวอักษรตรงส่วนใดส่วนหนึ่ง เช่น 'P' จะเจอ 'phone'
	rows, err := pdb.db.QueryContext(ctx, query, "%"+searchQuery+"%", storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (pdb *PostgresDatabase) GetAllProductsByStore(ctx context.Context, storeID int, sortOrder string) ([]Product, error) {
	// เรียงลำดับผลตามราคาตามค่าที่ส่งเข้ามา
	var orderByClause string
	if sortOrder == "asc" {
		orderByClause = "ORDER BY price ASC" // เรียงจากราคาน้อยไปมาก
	} else if sortOrder == "desc" {
		orderByClause = "ORDER BY price DESC" // เรียงจากราคามากไปน้อย
	} else {
		// หากไม่มีการส่งพารามิเตอร์ หรือค่าที่ไม่ถูกต้องให้ใช้การเรียงลำดับตามค่าเริ่มต้น
		orderByClause = "ORDER BY price ASC"
	}

	// สร้าง query ที่รองรับการเรียงลำดับ
	query := fmt.Sprintf(`
        SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path
        FROM product_info 
        WHERE store_id = $1
        %s;
    `, orderByClause)

	rows, err := pdb.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (pdb *PostgresDatabase) GetProductsByCategoryAndStore(ctx context.Context, storeID int, category string) ([]Product, error) {
	query := `
        SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path
        FROM product_info 
        WHERE store_id = $1 AND category = $2
        ORDER BY created_at DESC
    `
	rows, err := pdb.db.QueryContext(ctx, query, storeID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (pdb *PostgresDatabase) GetALLProductsByCategory(ctx context.Context, category string) ([]Product, error) {
	// ดึงข้อมูลสินค้าทุกตัวที่ตรงกับหมวดหมู่ที่ระบุ
	query := `
        SELECT id, product_name, price, quantity, created_at, updated_at, category, brand, model, store_id, is_recommended, image_path
        FROM product_info 
        WHERE category = $1
        ORDER BY created_at DESC
    `

	rows, err := pdb.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (bs *BookStore) GetAllStoreInfo(ctx context.Context) ([]StoreInfo, error) {
	return bs.db.GetAllStoreInfo(ctx)
}

func (bs *BookStore) GetStoreInfoByID(ctx context.Context, id int) (StoreInfo, error) {
	return bs.db.GetStoreInfoByID(ctx, id)
}

func (bs *BookStore) GetProductsByStore(ctx context.Context, storeID int) ([]Product, error) {
	return bs.db.GetProductsByStore(ctx, storeID)
}

func (bs *BookStore) GetNewProductsByStore(ctx context.Context, storeID int) ([]Product, error) {
	return bs.db.GetNewProductsByStore(ctx, storeID)
}

func (bs *BookStore) SearchProducts(ctx context.Context, searchQuery string) ([]Product, error) {
	return bs.db.SearchProducts(ctx, searchQuery)
}

func (bs *BookStore) GetProduct(ctx context.Context, id int) (Product, error) {
	return bs.db.GetProduct(ctx, id)
}

func (bs *BookStore) SearchProductsByStore(ctx context.Context, searchQuery string, storeID int) ([]Product, error) {
	return bs.db.SearchProductsByStore(ctx, searchQuery, storeID)
}

func (bs *BookStore) GetAllProductsByStore(ctx context.Context, storeID int, sortOrder string) ([]Product, error) {
	return bs.db.GetAllProductsByStore(ctx, storeID, sortOrder)
}

func (bs *BookStore) GetProductsByCategoryAndStore(ctx context.Context, storeID int, category string) ([]Product, error) {
	return bs.db.GetProductsByCategoryAndStore(ctx, storeID, category)
}

func (bs *BookStore) GetALLProductsByCategory(ctx context.Context, category string) ([]Product, error) {
	return bs.db.GetALLProductsByCategory(ctx, category)
}

func (pdb *PostgresDatabase) AddToCart(ctx context.Context, storeID, productID, quantity int) error {
	// ตรวจสอบว่ามีสินค้านี้อยู่ในตะกร้าหรือไม่
	var existingQuantity int
	query := `SELECT quantity FROM cart WHERE store_id = $1 AND product_id = $2 AND status = 'in_cart'`
	err := pdb.db.QueryRowContext(ctx, query, storeID, productID).Scan(&existingQuantity)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing item in cart: %v", err)
	}

	// ถ้ามีสินค้านี้อยู่แล้ว ให้เพิ่มจำนวน
	if err == nil {
		updateQuery := `UPDATE cart SET quantity = quantity + $1 WHERE store_id = $2 AND product_id = $3 AND status = 'in_cart'`
		_, err = pdb.db.ExecContext(ctx, updateQuery, quantity, storeID, productID)
		if err != nil {
			return fmt.Errorf("failed to update quantity in cart: %v", err)
		}
	} else {
		// ถ้ายังไม่มีในตะกร้า ให้เพิ่มรายการใหม่
		insertQuery := `
            INSERT INTO cart (store_id, product_id, quantity, added_at, status)
            VALUES ($1, $2, $3, $4, $5)
        `
		_, err := pdb.db.ExecContext(ctx, insertQuery, storeID, productID, quantity, time.Now(), "in_cart")
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %v", err)
		}
	}

	return nil
}

func (bs *BookStore) AddToCart(ctx context.Context, storeID, productID, quantity int) error {
	return bs.db.AddToCart(ctx, storeID, productID, quantity)
}

func (pdb *PostgresDatabase) GetCartItemsByStore(ctx context.Context, storeID int) ([]Product, error) {
	query := `SELECT p.id, p.product_name, p.price, c.quantity, p.created_at, p.updated_at, p.category, p.brand, p.model, p.store_id, p.is_recommended, p.image_path
              FROM cart c
              JOIN product_info p ON c.product_id = p.id
              WHERE p.store_id = $1 AND c.status = 'in_cart'`

	rows, err := pdb.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category,
			&product.Brand,
			&product.Model,
			&product.StoreID,
			&product.IsRecommended,
			&product.ImagePath,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product data: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return products, nil
}

func (bs *BookStore) GetCartItemsByStore(ctx context.Context, storeID int) ([]Product, error) {
	return bs.db.GetCartItemsByStore(ctx, storeID)
}

// DeleteProductFromCart ลบสินค้าจากตะกร้าสินค้าตาม productID
func (pdb *PostgresDatabase) DeleteProductFromCart(ctx context.Context, storeID, productID int) error {
	// Query สำหรับลบสินค้าจากตะกร้า
	query := `DELETE FROM cart WHERE store_id = $1 AND product_id = $2 AND status = 'in_cart'`
	// เรียกใช้คำสั่งลบจากฐานข้อมูล
	result, err := pdb.db.ExecContext(ctx, query, storeID, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product from cart: %v", err)
	}

	// ตรวจสอบว่ามีการลบข้อมูลหรือไม่
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	// ถ้าไม่มีแถวที่ถูกลบแสดงว่าไม่พบสินค้าที่จะลบ
	if rowsAffected == 0 {
		return fmt.Errorf("product not found in cart")
	}

	return nil
}

func (bs *BookStore) DeleteProductFromCart(ctx context.Context, storeID, productID int) error {
	return bs.db.DeleteProductFromCart(ctx, storeID, productID)
}

// ฟังก์ชันการชำระเงินและย้ายข้อมูลจากตะกร้าไปยัง order_history
func (pdb *PostgresDatabase) Checkout(ctx context.Context, storeID int) error {
	// เริ่มต้นการทำ transaction
	tx, err := pdb.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// ขอโต๊ะข้อมูลจากตะกร้าของ store ที่ระบุ
	query := `SELECT id, product_id, quantity FROM cart WHERE store_id = $1 AND status = 'in_cart'`
	rows, err := tx.QueryContext(ctx, query, storeID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch cart items: %v", err)
	}
	defer rows.Close()

	// สำหรับทุกสินค้าที่อยู่ในตะกร้า, ย้ายข้อมูลไปยัง order_history
	for rows.Next() {
		var cartID, productID, quantity int
		if err := rows.Scan(&cartID, &productID, &quantity); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan cart item: %v", err)
		}

		// Insert ข้อมูลลงใน order_history
		insertOrderQuery := `
            INSERT INTO order_history (store_id, product_id, quantity, status)
            VALUES ($1, $2, $3, $4)
        `
		_, err := tx.ExecContext(ctx, insertOrderQuery, storeID, productID, quantity, "ordered")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert order into order_history: %v", err)
		}

		// ลบสินค้านั้นออกจากตะกร้า
		deleteQuery := `DELETE FROM cart WHERE id = $1`
		_, err = tx.ExecContext(ctx, deleteQuery, cartID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete cart item: %v", err)
		}
	}

	// เช็คว่าไม่มีข้อผิดพลาดใด ๆ ก่อนที่จะทำการ commit
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error occurred while iterating over cart rows: %v", err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (bs *BookStore) Checkout(ctx context.Context, storeID int) error {
	return bs.db.Checkout(ctx, storeID)
}

func (pdb *PostgresDatabase) CheckoutCart(ctx context.Context, storeID int) error {
	query := `
        UPDATE cart
        SET status = 'checked_out', checked_out_at = $1
        WHERE store_id = $2 AND status = 'in_cart'
    `
	_, err := pdb.db.ExecContext(ctx, query, time.Now(), storeID)
	if err != nil {
		return fmt.Errorf("failed to checkout cart: %v", err)
	}
	return nil
}

func (bs *BookStore) CheckoutCart(ctx context.Context, storeID int) error {
	return bs.db.CheckoutCart(ctx, storeID)
}
