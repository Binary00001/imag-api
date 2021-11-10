package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//
//STRUCT DEFINITIONS
//

type job struct {
	Part_Number string 	`josn:"Part_Number"`
	Run *string					`json:"Run"`
	PO string						`json:"PO"`
	Item *string				`json:"Item"`
	Queue_Days	*int		`json:"Queue_Days"`
	Customer string 		`json:"Customer"`
	Priority int				`json:"Priority"`
	Comments *string		`json:"Comments"`
	Due_Date	string		`json:"Due_Date"`
	Run_Qty	string			`json:"Run_Qty"`
	Work_Center	string	`json:"Work_Center"`
	WC_Num	string			`json:"WC_Num"`
	WC_Name	string			`json:"WC_Name"`
}

type jobList []job

type employee struct {
	Job_Count		int				`json:"Job_Count"`
	Employee 		string		`json:"Employee"`
}

type employeeList []employee


//
//
//
// func getQuery(tsql string) {
// 	var temp job
// 	var tempList jobList

// 	ctx := context.Background()

// 	err := db.PingContext(ctx)
// 	if err != nil {
// 		log.Fatal("Could not establish connection: ", err.Error())
// 	}

// 	req := fmt.Sprintf(tsql)

// 	rows, err := db.QueryContext(ctx, req)
// 	if err != nil {
// 		log.Fatal("Error executing query: ", err.Error())
// 	}

// 	for rows.Next() {
// 		err := rows.Scan(
// 			&temp.Part_Number,
// 			&temp.Run,
// 			&temp.PO,
// 			&temp.Item,
// 			&temp.Queue_Days,
// 			&temp.Customer,
// 			&temp.Priority,
// 			&temp.Comments,
// 			&temp.Due_Date,
// 			&temp.Run_Qty,
// 			&temp.Work_Center,
// 			&temp.WC_Num,
// 			&temp.WC_Name,
// 		)

// 		if err != nil {
// 			log.Fatal("Error getting data: ", err.Error())
// 		}

// 		defer rows.Close()

// 		tempList = append(tempList, temp)
// 	}
// 	json.NewEncoder(w).Encode(tempList)
// }

//
//PARTS PER EMPLOYEE/DAY
//
func getEmployeeStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var temp employee
	var tempList employeeList

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT COUNT(OPREF) AS JOBS_COMPLETED, OPINSP AS EMPLOYEE FROM RnopTable INNER JOIN RunsTable ON RUNREF = OPREF AND RUNNO = OPRUN WHERE RUNPKPURGED = 0 AND OPCOMPDATE >= CAST(GETDATE() AS DATE) GROUP BY OPINSP")

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Uh oh: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.Job_Count,

			&temp.Employee)

		if err != nil {
			log.Fatal("Error getting info: ", err.Error())
		}

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)

}

//
//DEPARTMENT ROUTES
//

//
//GET TOP 20 PARTS IN DEPT BY DUE DATE
//
func getDept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]

	var temp job
	var tempList jobList

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT TOP 20 RTRIM(PART_NUMBER) AS PART_NUMBER, RTRIM(RUN) AS RUN, RTRIM(PO) AS PO, RTRIM(LTRIM(ITEM)) AS ITEM, DAYS_IN_QUEUE, RTRIM(CUSTOMER) AS CUSTOMER, PRIORITY, COMMENTS, CAST(CUST_REQ_DATE AS DATETIME) AS CUST_REQ_DATE, RUN_QTY, RTRIM(WORK_CENTER) AS WORK_CENTER, RTRIM(WC) AS WC, RTRIM(t2.WCNDESC) AS WC_NAME FROM QueueInfo INNER JOIN WcntTable AS t2 ON WC = t2.WCNNUM WHERE WC = '%s' ORDER BY CUST_REQ_DATE ASC", dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error executing query: ", err.Error())
	}

	

	for rows.Next() {
		err := rows.Scan(
			&temp.Part_Number,
			&temp.Run,
			&temp.PO,
			&temp.Item,
			&temp.Queue_Days,
			&temp.Customer,
			&temp.Priority,
			&temp.Comments,
			&temp.Due_Date,
			&temp.Run_Qty,
			&temp.Work_Center,
			&temp.WC_Num,
			&temp.WC_Name,
		)

		if err != nil {
			log.Fatal("Error: ", err.Error())
		}

		defer rows.Close()

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}

//
//DEPT BURNDOWN LIST
//
func getDeptBurndown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]
	var temp job
	var tempList jobList

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT RTRIM(PART_NUMBER) AS PART_NUMBER, RTRIM(RUN) AS RUN, RTRIM(PO) AS PO, RTRIM(LTRIM(ITEM)) AS ITEM, DAYS_IN_QUEUE, RTRIM(CUSTOMER) AS CUSTOMER, PRIORITY, COMMENTS, CAST(CUST_REQ_DATE AS DATETIME) AS CUST_REQ_DATE, RUN_QTY, RTRIM(WORK_CENTER) AS WORK_CENTER, RTRIM(WC) AS WC, RTRIM(t2.WCNDESC) AS WC_NAME FROM QueueInfo INNER JOIN WcntTable AS t2 ON WC = t2.WCNNUM WHERE SUBSTRING(COMMENTS, 1, 8) = 'BURNDOWN' AND WC = '%s' ORDER BY CUST_REQ_DATE ASC", dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error executing query: ", err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&temp.Part_Number,
			&temp.Run,
			&temp.PO,
			&temp.Item,
			&temp.Queue_Days,
			&temp.Customer,
			&temp.Priority,
			&temp.Comments,
			&temp.Due_Date,
			&temp.Run_Qty,
			&temp.Work_Center,
			&temp.WC_Num,
			&temp.WC_Name,
		)

		if err != nil {
			log.Fatal("Error: ", err.Error())
		}

		defer rows.Close()

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)

}


//
//BURNDOWN ALL
//
func getBurndown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var temp job
	var tempList jobList

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT RTRIM(PART_NUMBER) AS PART_NUMBER, RTRIM(RUN) AS RUN, RTRIM(PO) AS PO, RTRIM(LTRIM(ITEM)) AS ITEM, DAYS_IN_QUEUE, RTRIM(CUSTOMER) AS CUSTOMER, PRIORITY, COMMENTS, CAST(CUST_REQ_DATE AS DATETIME) AS CUST_REQ_DATE, RUN_QTY, RTRIM(WORK_CENTER) AS WORK_CENTER, RTRIM(WC) AS WC, RTRIM(t2.WCNDESC) AS WC_NAME FROM QueueInfo INNER JOIN WcntTable AS t2 ON WC = t2.WCNNUM WHERE SUBSTRING(COMMENTS, 1, 8) = 'BURNDOWN' ORDER BY WC ASC")

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error executing query: ", err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&temp.Part_Number,
			&temp.Run,
			&temp.PO,
			&temp.Item,
			&temp.Queue_Days,
			&temp.Customer,
			&temp.Priority,
			&temp.Comments,
			&temp.Due_Date,
			&temp.Run_Qty,
			&temp.Work_Center,
			&temp.WC_Num,
			&temp.WC_Name,
		)

		if err != nil {
			log.Fatal("Error: ", err.Error())
		}

		defer rows.Close()

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}