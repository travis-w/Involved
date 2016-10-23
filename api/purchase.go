package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
)

type Purchase struct{
	UserId int 			`json:"-"`
	ItemId int 			`json:"-"`
	Requested time.Time `json:"requested"`
	Fulfilled bool		`json:"fulfilled"`
	Details Item 		`json:"item"`
}

func purchaseRoute(w http.ResponseWriter, r *http.Request, user *User) {
	if user.CheckedInWith == 0 || user.Type != "host" {
		fmt.Fprintf(w, "[]")
		return
	}

	rows, err := db.Query("SELECT * FROM purchase p, item i WHERE p.item_id=i.item_id and p.user_id=?", user.Id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`, err)
		return
	}

	defer rows.Close()

	var purchases []*Purchase
	for rows.Next() {
		var p Purchase
		rows.Scan(
			&p.UserId,
			&p.ItemId,
			&p.Requested,
			&p.Fulfilled,
			&p.Details.Id,
			&p.Details.Cost,
			&p.Details.Seller,
			&p.Details.Description,
		)

		purchases = append(purchases, &p)
	}

	if rows.Err() != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`)
		return
	}

	jsonPurchase, err := json.Marshal(purchases)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "something went wrong"}`)
		return
	}

	fmt.Fprint(w, string(jsonPurchase))
}