package api

import (
	"fmt"
	"net/http"
	"time"
)

func logoutRoute(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: "",
		Expires: time.Now().Add(-1 * time.Second),
		HttpOnly: true,
	})

	fmt.Fprintf(w, "{msg: \"success\"}")
}