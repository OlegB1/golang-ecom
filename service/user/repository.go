package user

import (
	"database/sql"
	"fmt"

	"github.com/OlegB1/ecom/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var u *types.User
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u == nil {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Repository) CreateUser(user types.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (first_name,last_name,email,password) VALUES($1,$2,$3,$4)",
		user.FirstName, user.LastName, user.Email, user.Password,
	)
	return err
}

func (s *Repository) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var u *types.User
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u == nil {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}
