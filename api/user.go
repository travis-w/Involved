package api

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
)

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

func userRoute(w http.ResponseWriter, r *http.Request) {
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

	err = db.QueryRow("select * from user where id = ?", id).Scan(
		&parsedUser.Id,
		&parsedUser.Name,
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
