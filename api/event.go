package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"strings"
	"strconv"
	"errors"
)

type Event struct {
	Id 			int
	Creator		int
	Slots		int
	Divisions	int
	Description	string
	Start		time.Time
	EventType	string
	Needs		map[string]struct{}
	MetaData 	map[string]string
}

func (e *Event) MarshalJSON() ([]byte, error) {
	event := `{"id":%d,"creator":%d,"slots":%d,"divisions":%d,"type":"%s","description":"%s","start":"%s","needs":%s,"meta":%s}`
	needs := "[ "
	for key, _ := range e.Needs {
		needs += `"` + key + `",`
	}
	needs = needs[:len(needs)-1] + "]"

	meta := "{ "
	for key, value := range e.MetaData {
		meta += `"` + key + `":"` + value + `",`;
	}
	meta = meta[:len(meta)-1] + "}"

	return []byte(
		fmt.Sprintf(
			event,
			e.Id,
			e.Creator,
			e.Slots,
			e.Divisions,
			e.EventType,
			e.Description,
			e.Start.String(),
			needs,
			meta,
		),
	), nil
}

func eventRoute(w http.ResponseWriter, r *http.Request, user *User) {
	switch r.Method {
	case "POST":
		if user.Type == "seeker" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{"error": "You do not have permission to create events"}`)
			return
		}
		
		if user.CheckedInWith == 0 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{"error": "You must be verified by a center to create events"}`)
			return
		}

		event, err := validateEvent(r.URL.Query());

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "%v"}`, err)
			return
		}

		id, err := createEvent(user, event)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "%v"}`, err)
		}

		fmt.Fprintf(w, "{\"id\": %d}", id)
	case "GET":
		strId, ok := r.URL.Query()["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "no event id provided"}`)
			return
		}

		id, err := strconv.ParseInt(strId[0], 10, 0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "event id must be an integer"}`)
			return
		}

		event, err := getEvent(int(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			fmt.Fprintf(w, `{"error": "no such event found"}`)
			return
		}

		jsonEvent, err := json.Marshal(event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			fmt.Fprintf(w, `{"error": "error parsing event found"}`)
			return
		}

		fmt.Fprintf(w, string(jsonEvent))
	}
}

func validateEvent(params map[string][]string) (*Event, error) {
	require := []string{"type", "slots", "divisions", "description"}

	for _, key := range require {
		if _, ok := params[key]; !ok {
			return nil, errors.New("Required field " + key + " is missing")
		}
	}

	eventType, _ := params["type"]
	eventSlots, _ := params["slots"]
	eventDivs, _ := params["divisions"]
	eventDesc, _ := params["description"]
	eventNeeds, ok := params["needs"]
	needs := make(map[string]struct{})
	if ok {
		for _, val := range strings.Split(eventNeeds[0], ",") {
			needs[val] = struct{}{}
		}
	}

	fmt.Println(needs)
	meta, ok := params["meta"]
	metaData := make(map[string]string)

	if ok {
		for _, pair := range strings.Split(meta[0], "|") {
			key := strings.Split(pair, "===")[0]
			value := strings.Split(pair, "===")[1]

			metaData[key] = value
		}
	}

	slots, _ := strconv.ParseInt(eventSlots[0], 10, 0)
	divs, _ := strconv.ParseInt(eventDivs[0], 10, 0)

	return &Event{
		EventType: eventType[0],
		Slots: int(slots),
		Divisions: int(divs),
		Description: eventDesc[0],
		Needs: needs,
		MetaData: metaData,
		Start: time.Now(),
	}, nil
}

func createEvent(user *User, event *Event) (int, error) {
	res, err := db.Exec(
		"INSERT INTO event (user_id, availableSlots, maximumDivisions, description, type, created) VALUES (?, ?, ?, ?, ?, ?)",
		user.Id,
		event.Slots,
		event.Divisions,
		event.Description,
		event.EventType,
		event.Start,
	)

	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	for key, value := range event.MetaData {
		_, ok := event.Needs[key]
		_, err := db.Exec(
			"INSERT INTO event_meta (event_id, meta_key, value, isNeed) VALUES (?, ?, ?, ?)",
			int(id),
			key,
			value,
			ok,
		)

		if err != nil {
			fmt.Println(err)
		}
	}

	return int(id), err
}

func getEvent(id int) (*Event, error) {
	var event Event
	err := db.QueryRow("SELECT * FROM event WHERE event_id=?", id).Scan(
		&event.Id,
		&event.Creator,
		&event.Slots,
		&event.Divisions,
		&event.Description,
		&event.Start,
		&event.EventType,
	)

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT meta_key, value, isNeed FROM event_meta WHERE event_id=?", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	event.MetaData = make(map[string]string)
	event.Needs = make(map[string]struct{})
	for rows.Next() {
		var key string
		var value string
		var need bool

		err := rows.Scan(&key, &value, &need)

		if err != nil {
			return nil, err
		}

		event.MetaData[key] = value
		if need {
			event.Needs[key] = struct{}{}
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &event, nil
}
