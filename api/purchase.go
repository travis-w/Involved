package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"strconv"
	"errors"
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

	switch r.Method {
	case "GET":
		purchases, err := getPurchases(user.Id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "something went wrong"}`, err)
			return
		}

		jsonPurchase, err := json.Marshal(purchases)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "something went wrong"}`)
			return
		}

		fmt.Fprint(w, string(jsonPurchase))
	case "POST":
		item, ok := r.URL.Query()["item"]

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "must provide item identifier"}`)
			return
		}

		itemId, _ := strconv.ParseInt(item[0], 10, 0)

		err := makePurchase(user.Id, int(itemId))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "something went wrong"}`)
			return
		}

		fmt.Fprintf(w, `{"msg": "success"}`)
	case "PUT":
		if user.Type == "host" || user.Type == "seeker" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "you are not authorized to confirm purchases"}`)
			return
		}

		itemStr, iok := r.URL.Query()["item"]
		buyerStr, bok := r.URL.Query()["buyer"]

		if !iok || !bok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "must provide both a item identifier and buyer identifier"}`)
			return
		}

		item, _ := strconv.ParseInt(itemStr[0], 10, 0)
		buyer, _ := strconv.ParseInt(buyerStr[0], 10, 0)

		err := confirmPurchase(int(buyer), user.Id, int(item))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "could not confirm purchase"}`, err)
			return
		}

		fmt.Fprintf(w, `{"msg": "success"}`)
	}
}

func getPurchases(id int) ([]*Purchase, error) {
	rows, err := db.Query("SELECT * FROM purchase p, item i WHERE p.item_id=i.item_id and p.user_id=?", id)

	if err != nil {
		return nil, err
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
		return nil, err
	}

	return purchases, nil
}

func makePurchase(user, item int) error {
	_, err := db.Exec("INSERT INTO purchase (user_id, item_id) VALUES (?, ?)", user, item)

	return err
}

func confirmPurchase(buyer, seller, item int) error {
	res, err := db.Exec("UPDATE purchase p, item i SET fulfilled=1 " +
		"WHERE p.user_id=? and p.item_id=? and p.item_id=i.item_id and i.seller=?",
		buyer,
		item,
		seller,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil || rows == 0 {
		return errors.New("could not confirm purchase")
	}

	return nil
}