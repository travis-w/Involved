package api

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"strconv"
)

func queryEventsRoute(w http.ResponseWriter, r *http.Request, user *User) {
	switch r.Method {
	case "GET":
		events, err := getEvents(r.URL.Query())

		if err != nil {
			fmt.Fprintf(w, `{"error": "a problem occured with the query"}`)
			return
		}

		var asJSON []byte
		if events != nil {
			asJSON, err = json.Marshal(events)

			if err != nil {
				fmt.Fprintf(w, `{"error": "a problem occured during parsing"}`)
				return
			}
		}else{
			asJSON = []byte("[]")
		}

		fmt.Fprintf(w, string(asJSON))
	default:
		fmt.Fprintf(w, `{"error": "cannot access this route with that method"}`)
	}
}

func getEvents(params map[string][]string) ([]*Event, error) {
	queryStr := "SELECT e.* FROM event e "

	lat, doLocation := params["latitude"]
	if doLocation {
		lng, ok := params["longitude"]

		if !ok {
			return nil, errors.New("need both latitude and longitude to search by location")
		}

		rad, ok := params["radius"]

		if !ok {
			rad = []string{"10"}
		}

		realLat, _ := strconv.ParseFloat(lat[0], 64)
		realLng, _ := strconv.ParseFloat(lng[0], 64)
		realRad, _ := strconv.ParseFloat(rad[0], 64)

		realRad /= 50

		queryStr += ", event_location l WHERE e.event_id=l.event_id and "
		formula := "(l.latitude-%f)*(l.latitude-%f)+(l.longitude-%f)*(l.longitude-%f)<%f "
		queryStr += fmt.Sprintf(formula, realLat, realLat, realLng, realLng, realRad*realRad)

	} else {
		queryStr += "WHERE "
	}

	subQry_1 := "SELECT COALESCE(SUM(s.accepted),0) FROM seeker_event_response s, event e WHERE s.event_id=e.event_id"
	subQry_2 := "SELECT COALESCE(COUNT(*),0) FROM seeker_dependent d, event e, seeker_event_response s " +
				"WHERE d.user_id=s.user_id and e.event_id=s.event_id and s.accepted=1 "

	queryStr += "e.created > date_sub(current_timestamp, interval 1 day) "
	queryStr += "and e.maximumDivisions > (" + subQry_1 + ") or e.maximumDivisions=-1 "
	queryStr += "and e.availableSlots > (" + subQry_2 + ") or e.availableSlots=-1 "


	queryStr += "order by e.created"

	rows, err := db.Query(queryStr)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		rows.Scan(
			&event.Id,
			&event.Creator,
			&event.Slots,
			&event.Divisions,
			&event.Description,
			&event.Start,
			&event.EventType,
		)

		events = append(events, &event)
	}

	return events, nil
}