package paymentrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"
)

type PaymentRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) AddPayment(payment domains.Payment) error {
	fn := "paymentRepository"

	stmt := `
		INSERT INTO payments(order_id, payment, date)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(stmt, payment.OrderID, payment.Payment, payment.Date)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
