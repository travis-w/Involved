package api

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
	"errors"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
)

var salt string = "1m@#a$#bb&ikIcK@$$"

type User struct {
	Id    int			`json:"id"`
	Name  string		`json:"name"`
	Type  string		`json:"type"`
	Pic	  string		`json:"pic_url"`
	Desc  string		`json:"description"`
	Email string		`json:"email"`
	EmailVerified int	`json:"email_verified"`
	CheckedInWith int 	`json:"verified_by"`
	BelongsTo	  int 	`json:"member_of"`
}

func getUser(w http.ResponseWriter, r *http.Request) {
	strId, ok := r.URL.Query()["id"]

	if !ok {
		fmt.Fprintf(w, "{error: \"no user id provided\"}")
		return
	}

	id, err := strconv.ParseInt(strId[0], 10, 0)

	if err != nil {
		fmt.Fprintf(w, "{error: \"user id must be integer\"}")
		return
	}

	var parsedUser User
	var passHash string

	err = db.QueryRow("select * from user where id = ?", id).Scan(
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

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "{error: \"could not read from Database\"}")
		return
	}

	res, err := json.Marshal(parsedUser)

	if err != nil {
		fmt.Fprintf(w, "{error: \"could not parse user\"}")
		return
	}

	fmt.Fprintf(w, string(res))
}

func createUser(params map[string][]string) (int, error) {
	require := []string{"name", "email", "pass", "type"}

	for _, key := range require {
		if _, ok := params[key]; !ok {
			return -1, errors.New("Required field " + key + " is missing")
		}
	}

	name, _ := params["name"]
	email, _ := params["email"]
	pass, _ := params["pass"]
	userType, _ := params["type"]

	hash, err := scrypt.Key([]byte(pass[0]), []byte(salt), 16384, 8, 1, 32)

	if err != nil {
		return -1, errors.New("Failed to encrypt passphrase")
	}

	res, err := db.Exec(
		"INSERT INTO user (name, email, password, type) VALUES (?,?,?,?)",
		name[0],
		email[0],
		hex.EncodeToString(hash),
		userType[0])

	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func userRoute(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err == nil {
		_, err = userFromToken(cookie.Value)
	}

	switch r.Method {
	case "GET":
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "{error: \"must be logged in to view other users\"}")
			return
		}

		getUser(w, r)
	case "POST":
		id, err := createUser(r.URL.Query())

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "{error: \"failed to create user: %v\"}", err)
			return
		}

		fmt.Fprintf(w, "{id: %d}", id)
	case "PUT":
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "{error: \"must be logged in to view other users\"}")
			return
		}

		fmt.Fprintf(w, "{error: \"Updates not supported yet\"}")
	case "DELETE":
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "{error: \"must be logged in to view other users\"}")
			return
		}

		fmt.Fprintf(w, "{error: \"Updates not supported yet\"}")
	}
}
