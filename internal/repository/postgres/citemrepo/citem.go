package citemrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"

	"github.com/lib/pq"
)

var (
	ErrAlredyExists = fmt.Errorf("cart item alredy exists")
	ErrNotFound     = fmt.Errorf("cart item not found")
)

type CartItemRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *CartItemRepository {
	return &CartItemRepository{
		db: db,
	}
}

func (r *CartItemRepository) AddCartItem(citem domains.CartItem) error {
	fn := "cartItemRepository.AddCartItem"

	stmt := `
		INSERT INTO cart_items(user_id, product_id, quantity)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(stmt)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23505") {
			return fmt.Errorf("%s: %w", fn, ErrAlredyExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *CartItemRepository) GetCartItemsByUserID(userID uint32) ([]*domains.CartItem, error) {
	fn := "cartItemRepository.GetCartItemsByUserID"

	stmt := `
		SELECT user_id, product_id, quantity
		FROM cart_items
		WHERE user_id=$1
	`

	res, err := r.db.Query(stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var cartItems []*domains.CartItem
	for res.Next() {
		citem := &domains.CartItem{}
		err := res.Scan(&citem.UserID, &citem.ProductID, &citem.Quantity)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		cartItems = append(cartItems, citem)
	}

	return cartItems, nil
}

func (r *CartItemRepository) UpdateCartItemQuantity(userID, productID uint32, quantity int) error {
	fn := "cartItemRepository.UpdateCartItemQuantity"

	stmt := `
		UPDATE cart_items
		SET quantity=$1
		WHERE user_id=$2 AND product_id=$3
	`

	res, err := r.db.Exec(stmt, quantity, userID, productID)
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

func (r *CartItemRepository) AddCartItemQuantity(userID, productID uint32, quantity int) error {
	fn := "cartItemRepository.AddCartItemQuantity"

	stmt := `
		UPDATE cart_items
		SET quantity=quantity+$1
		WHERE user_id=$2 AND product_id=$3
	`

	res, err := r.db.Exec(stmt, quantity, userID, productID)
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

func (r *CartItemRepository) DeleteCartItemsByUserID(userID uint32) error {
	fn := "cartItemRepository.DeleteCartItemsByUserID"

	stmt := `
		DELETE FROM cart_items
		WHERE user_id=$1
	`

	_, err := r.db.Exec(stmt, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
