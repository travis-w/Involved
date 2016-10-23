package api

import (
	"fmt"
	"net/http"
)

func pointsRoute(w http.ResponseWriter, r *http.Request, user *User) {
	if user.CheckedInWith == 0 || user.Type != "host" {
		fmt.Fprint(w, `{"earned":0,"spent":0}`)
		return
	}

	var earned int
	var spent int

	err := db.QueryRow(
		"SELECT SUM(i.cost) FROM purchase p, item i " +
		"WHERE p.item_id=i.item_id and p.user_id=?", user.Id,
	).Scan(&spent)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`, err)
		return
	}

	err = db.QueryRow(
		"SELECT COUNT(*)*1000 FROM event WHERE user_id=?", user.Id,
	).Scan(&earned)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`, err)
		return
	}

	fmt.Fprintf(w, `{"earned": %d,"spent": %d}`, earned, spent);
}