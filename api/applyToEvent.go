package api

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
)

func applyToEventRoute(w http.ResponseWriter, r *http.Request, user *User) {
	if user.Type != "seeker" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"error": "You do not have permission to request invites to events"}`)
		return
	}

	if user.CheckedInWith == 0 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"error": "Your account must be verified before requesting invites"}`)
		return
	}

	params := r.URL.Query()
	eventStr, ok := params["event"]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"error": "Must provide event identifier"}`)
		return
	}

	attendStr, ok := params["attendees"]
	if !ok {
		attendStr = []string{"1"}
	}

	eventId, _ := strconv.ParseInt(eventStr[0], 10, 0)
	attendees, _ := strconv.ParseInt(attendStr[0], 10, 0)

	err := applyToEvent(user.Id, int(eventId), int(attendees))

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"error": "Something went wrong"}`)
		return
	}

	fmt.Fprintf(w, `{"msg": "success"}`)
}

func applyToEvent(user, event, attending int) error {
	res, err := db.Exec("INSERT INTO seeker_event_response VALUES (?,?,?,?)", event, user, attending, 0)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil || rows == 0 {
		return errors.New("Something went wrong")
	}

	return nil
}