package api

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"errors"
	"golang.org/x/crypto/scrypt"
	"time"
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

	if string(hash)[:len(storedHash)] != storedHash {
		return "", nil
	}

	token := generateToken()

	if err := storeToken(token, id); err != nil {
		return "", err
	}

	return token, nil
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
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

		fmt.Fprintf(w, "{msg: \"success\"}")
	}
}