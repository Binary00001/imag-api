package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *sql.DB


// THIS IS JUST FOR TESTING CONNECTION
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env", err.Error())
	}

	var server = os.Getenv("HOST")
	var port = 1433
	var user = os.Getenv("USER")
	var password = os.Getenv("PASSWORD")
	var database = os.Getenv("DATABASE")

	r := mux.NewRouter().StrictSlash(true)
	handler := cors.AllowAll().Handler(r)
	srv := &http.Server{
		Addr:    os.Getenv("ADDRESS"),
		Handler: handler,
	}

	// ROUTES
	r.HandleFunc("/", home)
	r.HandleFunc("/api/employee/stats", getEmployeeStats)
	r.HandleFunc("/api/dept/{dept}", getDept)
	r.HandleFunc("/api/dept/burndown/{dept}", getDeptBurndown)
	r.HandleFunc("/api/burndown", getBurndown)

	//PCM_ROUTES
	r.HandleFunc("/api/pcm", getPCMList)
	r.HandleFunc("/api/pcm/loc/{pcmLoc}", getPcmByLoc)

	//INVETORY 
	r.HandleFunc("/api/inv/part/{partNum}", lotsByPartNumber)
	r.HandleFunc("/api/inv/available", availableShip)

	// SETUP DATABASE
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)

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
