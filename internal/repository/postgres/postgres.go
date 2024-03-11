package postgres

import (
	"database/sql"
	"fmt"
	"shop/internal/config"
	"shop/internal/repository/postgres/categoryrepo"
	"shop/internal/repository/postgres/oitemrepo"
	"shop/internal/repository/postgres/orderrepo"
	"shop/internal/repository/postgres/paymentrepo"
	"shop/internal/repository/postgres/productrepo"
	"shop/internal/repository/postgres/userrepo"

	_ "github.com/lib/pq"
)

type Repository struct {
	*categoryrepo.ProductCategoryRepository
	*productrepo.ProductRepository
	*userrepo.UserRepository
	*orderrepo.OrderRepository
	*oitemrepo.OrderItemRepository
	*paymentrepo.PaymentRepository
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
		userrepo.New(db),
		orderrepo.New(db),
		oitemrepo.New(db),
		paymentrepo.New(db),
	}, nil
}
