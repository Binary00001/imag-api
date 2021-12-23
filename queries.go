package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getJobs(tsql string) []Job {

	var temp Job
	var tempList []Job

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprint(tsql)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.PartRef,
			&temp.PartNum,
			&temp.Run,
			&temp.Quantity,
			&temp.Comments,
			&temp.Customer,
			&temp.PO,
			&temp.Item,
			&temp.CustDate,
			&temp.WCNum,
			&temp.WCName,
			&temp.Priority,
			&temp.QueueDiff,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	return tempList
}


func getDeptStats(w http.ResponseWriter, r *http.Request) {

	dept := mux.Vars(r)["dept"]
	w.Header().Set("Content-Type", "application/json")

	var temp DeptStats


	temp.JobCount = completedJobs(dept)

	temp.PartCount = completedParts(dept)

	temp.Goal = dailyGoal(dept)

	
	json.NewEncoder(w).Encode(temp)
	}


