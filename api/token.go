package api

import (
    "math/rand"
    "time"
)

var hasSeed bool = false

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateToken() string {
	if !hasSeed {
		rand.Seed(time.Now().UTC().UnixNano())
		hasSeed = true
	}
	b := make([]rune, 128)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func storeToken(token string, user int) error {
	_, err := db.Exec("INSERT INTO token (user_id, value) VALUES (?,?)", user, token)
	return err
}

func userFromToken(token string) (*User, error) {
	var parsedUser User
	var passHash string

	err := db.QueryRow("select user.* from user, token where value=? and user.id=user_id", token).Scan(
		&parsedUser.Id,
		&parsedUser.Name,
		&passHash,
		&parsedUser.Email,
		&parsedUser.Pic,
		&parsedUser.EmailVerified,
		&parsedUser.CheckedInWith,
		&parsedUser.BelongsTo,
		&parsedUser.Desc,
		&parsedUser.Type)

	return &parsedUser, err
}