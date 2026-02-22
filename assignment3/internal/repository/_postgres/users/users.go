package users

import (
	"errors"

	"golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db *_postgres.Dialect
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = []modules.User{}
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *modules.User) (int, error) {
	query := `
	INSERT INTO users (name, email, age)
	VALUES ($1, $2, $3)
	RETURNING id`

	var id int
	err := r.db.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Age,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(user *modules.User) error {
	result, err := r.db.DB.Exec(
		`UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4`,
		user.Name,
		user.Email,
		user.Age,
		user.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return 0, errors.New("user not found")
	}

	return rows, nil
}
