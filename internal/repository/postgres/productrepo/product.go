package productrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"

	"github.com/lib/pq"
)

var (
	ErrNotFound   = fmt.Errorf("product not fount")
	ErrNoCategory = fmt.Errorf("no product category")
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

func (r *ProductRepository) GetProductByID(id uint32) (*domains.Product, error) {
	fn := "productRepository.GetProductByID"

	stmt := `
		SELECT id, name, description, price, quantiry, category_id, image_path
		FROM products
		WHERE id=$1
	`

	product := &domains.Product{}
	row := r.db.QueryRow(stmt, id)
	err := row.Scan(&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.CategoryID,
		&product.ImagePath,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return product, nil
}

func (r *ProductRepository) GetAllProducts() ([]*domains.Product, error) {
	fn := "productRepository.GetAllProducts"

	stmt := `
		SELECT id, name, description, price, quantiry, category_id, image_path
		FROM products 
	`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var products []*domains.Product
	for rows.Next() {
		product := &domains.Product{}
		err := rows.Scan(&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CategoryID,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductsByCategoryID(categoryID uint32) ([]*domains.Product, error) {
	fn := "productRepository.GetProductsByCategory"

	stmt := `
		SELECT id, name, description, price, quantiry, category_id, image_path
		FROM products
		WHERE category_id=$1
	`

	rows, err := r.db.Query(stmt, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var products []*domains.Product
	for rows.Next() {
		product := &domains.Product{}
		err := rows.Scan(&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CategoryID,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) updateField(id uint32, field string, value any) error {
	stmt := fmt.Sprintf(`
		UPDATE products
		SET %s=$1
		WHERE id=$2
	`, field)

	res, err := r.db.Exec(stmt, field, value)
	if err != nil {
		return err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *ProductRepository) UpdateProductName(id uint32, name string) error {
	fn := "productRepository.UpdateProductName"
	if err := r.updateField(id, "name", name); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductDescription(id uint32, description string) error {
	fn := "productRepository.UpdateProductDescription"
	if err := r.updateField(id, "description", description); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductPrice(id uint32, price float64) error {
	fn := "productRepository.UpdateProductPrice"
	if err := r.updateField(id, "price", price); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductQuantity(id uint32, quantity int) error {
	fn := "productRepository.UpdateProductQuantity"
	if err := r.updateField(id, "quantity", quantity); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *ProductRepository) AddProductQuantity(id uint32, quantity int) error {
	fn := "productRepository.AddProductQuantity"

	stmt := `
		UPDATE products
		SET quantity=quantity+$1
		WHERE id=$2
	`

	res, err := r.db.Exec(stmt, quantity, id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		return fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	return nil
}

func (r *ProductRepository) UpdateProductCategory(id, categoryID uint32) error {
	fn := "productRepository.UpdateProductCategory"
	err := r.updateField(id, "category_id", categoryID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23503") {
			return fmt.Errorf("%s: %w", fn, ErrNoCategory)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductImagePath(id uint32, imagePath string) error {
	fn := "productRepository.UpdateProductImagePath"
	if err := r.updateField(id, "image_path", imagePath); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *ProductRepository) DeleteProductByID(id uint32) (bool, error) {
	fn := "productRepository.DeleteProductByID"

	stmt := `
		DELETE FROM products
		WHERE id=$1
	`

	res, err := r.db.Exec(stmt, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		return false, fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	return true, nil
}
