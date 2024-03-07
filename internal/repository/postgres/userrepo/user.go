package userrepo

import (
	"database/sql"
	"fmt"
	"shop/internal/domains"

	"github.com/lib/pq"
)

var (
	ErrAlredyExists = fmt.Errorf("user alredy exists")
	ErrNotFound     = fmt.Errorf("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddUser(user domains.User) error {
	fn := "userRepository.AddUser"

	stmt := `
		INSERT INTO users(first_name, last_name, email, password, phone)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(stmt, user.FirstName, user.LastName, user.Email, user.Password, user.Phone)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23505") {
			return fmt.Errorf("%s: %w", fn, ErrAlredyExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *UserRepository) GetUserByID(id uint32) (*domains.User, error) {
	fn := "userRepository.GetUserByID"

	stmt := `
		SELECT id, first_name, last_name, email, password, phone
		FROM users
		WHERE id=$1
	`

	user := &domains.User{}
	row := r.db.QueryRow(stmt, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domains.User, error) {
	fn := "userRepository.GetUserByEmail"

	stmt := `
		SELECT id, first_name, last_name, email, password, phone
		FROM users
		WHERE email=$1
	`

	user := &domains.User{}
	row := r.db.QueryRow(stmt, email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (r *UserRepository) updateField(id uint32, field string, value any) error {
	stmt := fmt.Sprintf(`
		UPDATE users
		SET %s=$1
		WHERE id=$2
	`, field)

	res, err := r.db.Exec(stmt, value, id)
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

func (r *UserRepository) UpdateUserFirstName(id uint32, firstName string) error {
	fn := "userRepository.UpdateUserFirstName"
	if err := r.updateField(id, "first_name", firstName); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *UserRepository) UpdateUserLastName(id uint32, lastName string) error {
	fn := "userRepository.UpdateUserLastName"
	if err := r.updateField(id, "last_name", lastName); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *UserRepository) UpdateUserPassword(id uint32, password string) error {
	fn := "userRepository.UpdateUserPassword"
	if err := r.updateField(id, "password", password); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *UserRepository) DeleteUserByID(id uint32) (bool, error) {
	fn := "userRepository.DeleteUserByID"

	stmt := `
		DELETE FROM users
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
