package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB
var server = "IMASQL01"
var port = 1433
var user = "esisa"
var password = "CNIAMI000"
var database = "Esi2000Db"

// THIS IS JUST FOR TESTING CONNECTION
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")

}


func main() {
	r := mux.NewRouter().StrictSlash(true)
	handler := cors.AllowAll().Handler(r)
	srv := &http.Server{
		Addr: "10.25.1.73:4004",
		Handler: handler,
	}

	// ROUTES
	r.HandleFunc("/", home)
	r.HandleFunc("/api/employee/stats", getEmployeeStats)
	r.HandleFunc("/api/dept/{dept}", getDept)
	r.HandleFunc("/api/dept/burndown/{dept}", getDeptBurndown)
	r.HandleFunc("/api/burndown", getBurndown)

	// SETUP DATABASE
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)

	var err error	

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not connect to database: ", err.Error())
	}
	fmt.Println("Connected")
	
	log.Fatal(srv.ListenAndServe())
}