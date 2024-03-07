package categoryrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"

	"github.com/lib/pq"
)

var (
	ErrAlredyExists = fmt.Errorf("product category alredy exists")
	ErrNotFound     = fmt.Errorf("product category not found")
)

type ProductCategoryRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *ProductCategoryRepository {
	return &ProductCategoryRepository{
		db: db,
	}
}

func (r *ProductCategoryRepository) AddCategory(category domains.ProductCategory) error {
	fn := "productCategoryRepository.AddCategory"

	stmt := `
		INSERT INTO product_categories(name, description)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(stmt, category.Name, category.Description)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23505") {
			return fmt.Errorf("%s: %w", fn, ErrAlredyExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ProductCategoryRepository) GetAllCategories() ([]*domains.ProductCategory, error) {
	fn := "procuctCategoryRepository.GetAll"

	stmt := `
		SELECT id, name, description FROM product_categories
	`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var categories []*domains.ProductCategory

	for rows.Next() {
		category := &domains.ProductCategory{}
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *ProductCategoryRepository) GetCategoryByID(id uint32) (*domains.ProductCategory, error) {
	fn := "productCategory.GetCategoryByID"

	stmt := `
		SELECT id, name, description FROM product_categories
		WHERE id=$1
	`

	category := &domains.ProductCategory{}
	row := r.db.QueryRow(stmt, id)
	err := row.Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return category, nil
}
