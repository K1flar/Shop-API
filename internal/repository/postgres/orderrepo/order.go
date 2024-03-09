package orderrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"

	"github.com/lib/pq"
)

var (
	ErrInvalidStatus       = fmt.Errorf("invalid status order")
	ErrInvalidPickupMethod = fmt.Errorf("invalid pickup method")
	ErrNotFound            = fmt.Errorf("order not found")
)

type OrderRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) AddOrder(order *domains.Order) error {
	fn := "orderRepository.AddOrder"

	stmt := `
		INSERT INTO orders(user_id, date, price, status, pickup_method, delivery_address) 
		VALUES($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(stmt,
		order.UserID,
		order.Date,
		order.Price,
		order.Status,
		order.PickupMethod,
		order.DeliveryAddress,
	)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			switch err.Constraint {
			case "orders_status_check":
				return fmt.Errorf("%s: %w", fn, ErrInvalidStatus)
			case "orders_pickup_method_check":
				return fmt.Errorf("%s: %w", fn, ErrInvalidPickupMethod)
			}
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *OrderRepository) GetOrderByID(id uint32) (*domains.Order, error) {
	fn := "orderRepository.GetOrderByID"

	stmt := `
		SELECT id, user_id, date, price, status, pickup_method, delivery_address
		FROM orders
		WHERE id=$1
	`

	order := &domains.Order{}
	row := r.db.QueryRow(stmt, id)
	err := row.Scan(&order.ID,
		&order.UserID,
		&order.Date,
		&order.Price,
		&order.Status,
		&order.PickupMethod,
		&order.DeliveryAddress,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID uint32) ([]*domains.Order, error) {
	fn := "orderRepository.GetOrdersByUserID"

	stmt := `
		SELECT id, user_id, date, price, status, pickup_method, delivery_address
		FROM orders
		WHERE user_id=$1
	`

	res, err := r.db.Query(stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var orders []*domains.Order
	for res.Next() {
		order := &domains.Order{}
		err := res.Scan(&order.ID,
			&order.UserID,
			&order.Date,
			&order.Price,
			&order.Status,
			&order.PickupMethod,
			&order.DeliveryAddress,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) getOrdersByUserIDWtighStatus(userID uint32, status string) ([]*domains.Order, error) {
	stmt := `
		SELECT id, user_id, date, price, status, pickup_method, delivery_address
		FROM orders
		WHERE user_id=$1 AND status=$2
	`

	res, err := r.db.Query(stmt, userID, status)
	if err != nil {
		return nil, err
	}

	var orders []*domains.Order
	for res.Next() {
		order := &domains.Order{}
		err := res.Scan(&order.ID,
			&order.UserID,
			&order.Date,
			&order.Price,
			&order.Status,
			&order.PickupMethod,
			&order.DeliveryAddress,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetPendingOrdersByUserID(userID uint32) ([]*domains.Order, error) {
	fn := "orderRepository.GetPendingOrdersByUserID"
	orders, err := r.getOrdersByUserIDWtighStatus(userID, "pending")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orders, nil
}

func (r *OrderRepository) GetPaidOrdersByUserID(userID uint32) ([]*domains.Order, error) {
	fn := "orderRepository.GetPaidOrdersByUserID"
	orders, err := r.getOrdersByUserIDWtighStatus(userID, "paid")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orders, nil
}

func (r *OrderRepository) GetCanceledOrdersByUserID(userID uint32) ([]*domains.Order, error) {
	fn := "orderRepository.GetCanceledOrdersByUserID"
	orders, err := r.getOrdersByUserIDWtighStatus(userID, "canceled")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(id uint32, status string) error {
	fn := "orderRepository.UpdateOrderStatus"

	stmt := `
		UPDATE orders
		SET status=$1
		WHERE id=$2
	`

	_, err := r.db.Exec(stmt, status, id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidStatus)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
