package postgres

import (
	"database/sql"
	"fmt"
	"shop/internal/config"
	"shop/internal/repository/postgres/categoryrepo"
	"shop/internal/repository/postgres/productrepo"

	_ "github.com/lib/pq"
)

type Repository struct {
	*categoryrepo.ProductCategoryRepository
	*productrepo.ProductRepository
}

func New(cfg config.DataBase) (*Repository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Repository{
		categoryrepo.New(db),
		productrepo.New(db),
	}, nil
}