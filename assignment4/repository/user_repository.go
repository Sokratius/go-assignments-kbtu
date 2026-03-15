package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"assignment4/models"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int, filters map[string]string, orderBy string) (models.PaginatedResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	allowedOrders := map[string]bool{"id": true, "name": true, "email": true, "gender": true, "birth_date": true}
	if orderBy == "" {
		orderBy = "id"
	}
	if !allowedOrders[orderBy] {
		orderBy = "id"
	}

	base := "FROM users"
	whereParts := []string{}
	args := []interface{}{}
	argIdx := 1

	for k, v := range filters {
		if v == "" {
			continue
		}
		switch k {
		case "id":
			if u, err := uuid.Parse(v); err == nil {
				whereParts = append(whereParts, fmt.Sprintf("id = $%d", argIdx))
				args = append(args, u)
				argIdx++
			}
		case "name":
			whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", argIdx))
			args = append(args, "%"+v+"%")
			argIdx++
		case "email":
			whereParts = append(whereParts, fmt.Sprintf("email ILIKE $%d", argIdx))
			args = append(args, "%"+v+"%")
			argIdx++
		case "gender":
			whereParts = append(whereParts, fmt.Sprintf("gender ILIKE $%d", argIdx))
			args = append(args, "%"+v+"%")
			argIdx++
		case "birth_date":
			whereParts = append(whereParts, fmt.Sprintf("birth_date = $%d", argIdx))
			args = append(args, v)
			argIdx++
		}
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = " WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + base + whereClause
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return models.PaginatedResponse{}, err
	}

	query := "SELECT id, name, email, gender, birth_date " + base + whereClause + fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", orderBy, argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return models.PaginatedResponse{}, err
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(userID1 uuid.UUID, userID2 uuid.UUID) ([]models.User, error) {
	query := `
SELECT u.id, u.name, u.email, u.gender, u.birth_date
FROM users u
JOIN user_friends f1 ON u.id = f1.friend_id
JOIN user_friends f2 ON u.id = f2.friend_id
WHERE f1.user_id = $1 AND f2.user_id = $2
`

	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
