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
	_, err := db.Exec("INSERT INTO token (id, value) VALUES (?,?)", user, token)
	return err
}

func userFromToken(token string) (int, error) {
	var id int = -1
	err := db.QueryRow("select id from token where value = ?", token).Scan(&id)
	return id, err
}