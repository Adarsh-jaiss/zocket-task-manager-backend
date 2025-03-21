package users

import (
	"database/sql"
	"errors"

	"github.com/adarsh-jaiss/zocket/types"
)

var (
	ErrEmailExists = errors.New("email already exists")
)

func GetUserByEmailAndPassword(db *sql.DB, email string) (types.SignInRequest, error) {
	var user types.SignInRequest
	query := `
		SELECT user_id, email, password
		FROM users WHERE email = $1
	`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return types.SignInRequest{}, err
	}
	return user, nil
}

func GetUserFromStore(db *sql.DB, userID int) (types.User, error) {
	var user types.User
	query := `
		SELECT user_id, email, first_name, last_name, created_at, logged_in_at FROM users WHERE user_id = $1
	`
	err := db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.LoggedInAt,
	)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}

func CreateUserInStore(db *sql.DB, user types.User) (int, error) {
	// Check if email already exists
	_, err := GetUserByEmailAndPassword(db, user.Email)
	if err == nil {
		return 0, ErrEmailExists
	} else if err != sql.ErrNoRows {
		return 0, err
	}
	query := `
		INSERT INTO users (email, password, first_name, last_name, created_at, logged_in_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING user_id
	`
	var userID int
	err = db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}
	return userID, nil
}

func GetAllUsers(db *sql.DB) ([]types.User, error) {
	var users []types.User
	query := `
		SELECT user_id, email, first_name, last_name
		FROM users
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
