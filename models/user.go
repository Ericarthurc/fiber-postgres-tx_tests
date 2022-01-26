package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Userseats   int       `json:"userseats"`
	Servicetime time.Time `json:"servicetime"`
	Serviceid   uuid.UUID `json:"serviceid"`
}

// func GetUsers() {}

// func GetUserByID(db *pgxpool.Pool, id string) (*User, error) {
// 	var user User
// 	err := pgxscan.Get(context.Background(), db, &user, "SELECT * FROM users WHERE id = $1", id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func GetUserByUsername(db *pgxpool.Pool, username string) (*User, error) {
// 	var user User
// 	err := pgxscan.Get(context.Background(), db, &user, "SELECT * FROM users WHERE username = $1", username)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func CheckForToken(db *pgxpool.Pool, id string, token string) bool {
// 	var tokens []string
// 	err := pgxscan.Get(context.Background(), db, &tokens, "SELECT tokens FROM users WHERE id = $1", &id)
// 	if err != nil {
// 		return false
// 	}

// 	for _, v := range tokens {
// 		if v == token {
// 			return true
// 		}
// 	}
// 	return false
// }

// func LoginUser() {}

// func LogoutUser(db *pgxpool.Pool, id string, token string) error {
// 	_, err := db.Exec(context.Background(), "UPDATE users SET tokens = array_remove(tokens, $2) WHERE id = $1", &id, &token)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func CreateUser(db *pgxpool.Pool, u User) (*User, error) {
// 	var user User
// 	err := pgxscan.Get(context.Background(), db, &user, "INSERT INTO users (id, username, password, role, tokens) VALUES ($1, $2, $3, $4) RETURNING *", &u.ID, &u.Username, &u.Password, &u.Role)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // func UpdateUserTokens(db *pgxpool.Pool, u User) (*User, error) {
// // 	var user User
// // 	err := pgxscan.Get(context.Background(), db, &user, "UPDATE users SET tokens = $2 WHERE id = $1 RETURNING *", &u.ID)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return &user, nil
// // }

// func UpdateUser() {}

// func DeleteUser() {}
