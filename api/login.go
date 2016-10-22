package api

import (
	"fmt"
	"net/http"
	"errors"
	"golang.org/x/crypto/scrypt"
	"time"
	"encoding/hex"
	"encoding/json"
)

func login(params map[string][]string) (string, error) {
	require := []string{"email", "pass"}

	for _, key := range require {
		if _, ok := params[key]; !ok {
			return "", errors.New("Required field " + key + " is missing")
		}
	}

	email, _ := params["email"]
	pass, _ := params["pass"]

	hash, err := scrypt.Key([]byte(pass[0]), []byte(salt), 16384, 8, 1, 32)
	var storedHash string
	var id int

	err = db.QueryRow("select password, id from user where email=?", email[0]).Scan(&storedHash, &id)

	if err != nil {
		return "", err
	}

	if hex.EncodeToString(hash) != storedHash {
		return "", nil
	}

	token := generateToken()

	if err := storeToken(token, id); err != nil {
		return "", err
	}

	return token, nil
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err == nil {
		user, err := userFromToken(cookie.Value)

		if err == nil {
			res, err := json.Marshal(user)

			if err != nil {
				fmt.Fprintf(w, "{error: \"could not parse user\"}")
				return
			}

			fmt.Fprintf(w, string(res))
			return
		}
	}

	switch r.Method {
	case "POST":
		token, err := login(r.URL.Query())

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{error: \"login failed: %v\"}", err)
			return
		}

		if token == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{error: \"login failed: wrong email or password\"}")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value: token,
			Expires: time.Now().Add(180 * 24 * time.Hour),
			HttpOnly: true,
		})

		user, err := userFromToken(token)

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "{error: \"could not read user\"}")
			return
		}

		res, err := json.Marshal(user)

		if err != nil {
			fmt.Fprintf(w, "{error: \"could not parse user\"}")
			return
		}

		fmt.Fprintf(w, string(res))
	case "GET":
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{error: \"not logged in\"}")
	}
}