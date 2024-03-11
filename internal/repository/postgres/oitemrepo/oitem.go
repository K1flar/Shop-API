package oitemrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"
)

type OrderItemRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r *OrderItemRepository) AddOrderItem(oitem domains.OrderItem) error {
	fn := "orderItemRepository.AddOrderItem"

	stmt := `
		INSERT INTO order_items(order_id, product_id, quantity, product_price)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(stmt, oitem.OrderID, oitem.ProductID, oitem.Quantity, oitem.ProductPrice)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *OrderItemRepository) GetOrderItemsByOrderID(orderID uint32) ([]*domains.OrderItem, error) {
	fn := "orderItemRepository.GetOrderItemsByOrderID"

	stmt := `
		SELECT order_id, product_id, quantity, product_price
		FROM order_items
		WHERE order_id=$1
	`

	res, err := r.db.Query(stmt, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var oitems []*domains.OrderItem
	for res.Next() {
		oitem := &domains.OrderItem{}
		err := res.Scan(oitem.OrderID, oitem.ProductID, oitem.Quantity, oitem.ProductPrice)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		oitems = append(oitems, oitem)
	}

	return oitems, nil
}
