package productrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"
)

type ProductRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) AddProduct(product domains.Product) error {
	fn := "productRepository.AddProduct"

	stmt := `
		INSERT INTO products(name, description, price, quantity, category_id, image_path)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.db.Exec(stmt, product.Name, product.Description, product.Price, product.Quantity, product.CategoryID, product.ImagePath)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
