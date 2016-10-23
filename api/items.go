package api

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type Item struct{
	Id int  			`json:"id"`
	Cost int  			`json:"cost"`
	Seller int 			`json:"seller"`
	Description string	`json:"description"`
}

func itemsRoute(w http.ResponseWriter, r *http.Request, user *User) {
	if user.CheckedInWith == 0 || user.Type != "host" {
		fmt.Fprintf(w, "[]")
		return
	}

	rows, err := db.Query("SELECT * FROM item")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`, err)
		return
	}

	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var i Item
		rows.Scan(
			&i.Id,
			&i.Cost,
			&i.Seller,
			&i.Description,
		)

		items = append(items, &i)
	}

	if rows.Err() != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`)
		return
	}

	jsonItem, err := json.Marshal(items)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`)
		return
	}

	fmt.Fprint(w, string(jsonItem))
}