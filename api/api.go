package api

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var routes = make(map[string]func(http.ResponseWriter,*http.Request))

func RegisterRoute(route string, handler func(http.ResponseWriter,*http.Request)) {
	routes[route] = handler
}

func Init() {
	tempDB, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/involved?parseTime=true")
	if err != nil {
		fmt.Println(err)
	}

	db = tempDB

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
		return
	}

	RegisterRoute("login", loginRoute)
	RegisterRoute("logout", logoutRoute)
	RegisterRoute("user", userRoute)
	RegisterRoute("event", requireAuth(eventRoute))
	RegisterRoute("event/apply", requireAuth(applyToEventRoute))
	RegisterRoute("event/accept", requireAuth(acceptToEventRoute))
	RegisterRoute("events", requireAuth(queryEventsRoute))
	RegisterRoute("user/verify", requireAuth(verifyUserRoute))
	RegisterRoute("user/purchases", requireAuth(purchaseRoute))
	RegisterRoute("user/points", requireAuth(pointsRoute))
	RegisterRoute("items", requireAuth(itemsRoute))
}

func Shutdown() {
	db.Close()
}

func Request(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/"):]
	fn,ok := routes[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "%s is not a registered API route."}`, path)
		return
	}

	fn(w, r)
}
