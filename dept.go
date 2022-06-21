package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]

	var temp Part
	var tempList []Part

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprintf(`
			SELECT DISTINCT TOP 20
				RUNREF,
				RUNRTNUM, 
				RUNNO,
				RUNQTY,
				SOCUST,
				--SOPO, 
				--RASOITEM, 
				ISNULL(AGPMCOMMENTS, '') AGPMCOMMENTS,
				RUNPRIORITY, 
				OPSCHEDDATE,
				ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), 0) DTDIFF,
				WCNDESC
				FROM RunsTable
				INNER JOIN RnopTable ON OPREF=RUNREF AND OPRUN= RUNNO AND RUNOPCUR=OPNO
				INNER JOIN PartTable ON PARTREF=RUNREF 
				INNER JOIN RnalTable ON RUNREF=RAREF AND RUNNO=RARUN
				INNER JOIN SohdTable ON SONUMBER=RASO
				INNER JOIN WcntTable ON OPCENTER = WCNREF
				LEFT OUTER JOIN AgcmTable ON AGPART=RUNRTNUM AND AGRUN=RUNNO
				LEFT OUTER JOIN SoitTable ON ITPART=PARTREF AND ITSO=RASO

				WHERE OPCENTER LIKE '%s'
				AND OPCOMPLETE = 0
				AND SOTYPE != 'F'
				ORDER BY RUNPRIORITY, OPSCHEDDATE`, dept)
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
			&temp.Customer,
			&temp.Comments,
			&temp.Priority,
			&temp.SchedDate,
			&temp.QueueDiff,
			&temp.WCName,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}

func getQueueList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]

	var temp Part
	var tempList []Part

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprintf(`
			SELECT DISTINCT TOP 20
				RUNREF,
				RUNRTNUM, 
				RUNNO,
				RUNQTY,
				SOCUST,
				--SOPO, 
				--RASOITEM, 
				ISNULL(AGPMCOMMENTS, '') AGPMCOMMENTS,
				RUNPRIORITY, 
				OPSCHEDDATE,
				ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), 0) DTDIFF,
				WCNDESC
				FROM RunsTable
				INNER JOIN RnopTable ON OPREF=RUNREF AND OPRUN= RUNNO AND RUNOPCUR=OPNO
				INNER JOIN PartTable ON PARTREF=RUNREF 
				INNER JOIN RnalTable ON RUNREF=RAREF AND RUNNO=RARUN
				INNER JOIN SohdTable ON SONUMBER=RASO
				INNER JOIN WcntTable ON OPCENTER = WCNREF
				LEFT OUTER JOIN AgcmTable ON AGPART=RUNRTNUM AND AGRUN=RUNNO
				LEFT OUTER JOIN SoitTable ON ITPART=PARTREF AND ITSO=RASO

				WHERE OPCENTER LIKE '%s'
				AND OPCOMPLETE = 0
				ORDER BY RUNPRIORITY, OPSCHEDDATE`, dept)
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
			&temp.Customer,
			&temp.Comments,
			&temp.Priority,
			&temp.SchedDate,
			&temp.QueueDiff,
			&temp.WCName,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}

//DAILY DEPT STATISTICS (JOBS COMPLETED, PARTS COMPLETED, DAILY GOAL)
func getDeptStats(w http.ResponseWriter, r *http.Request) {

	dept := mux.Vars(r)["dept"]
	w.Header().Set("Content-Type", "application/json")

	var temp DeptStats

	temp.JobCount = completedJobs(dept)

	temp.PartCount = completedParts(dept)

	temp.Goal = dailyGoal(dept)

	json.NewEncoder(w).Encode(temp)
}

//

//DAILY JOB COUNT FROM THE PAST WEEK
func getChartData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dept := mux.Vars(r)["dept"]

	var temp ChartData
	var tempList []ChartData

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprintf(`
		SELECT COUNT(OPREF) AS JOB_COUNT, 
    CAST(OPCOMPDATE AS DATE) AS DATE 
		FROM RnopTable
    INNER JOIN RunsTable ON RUNREF = OPREF
    AND RUNNO = OPRUN
    WHERE RUNPKPURGED = 0
    AND OPCOMPDATE >= CAST(GETDATE() AS DATETIME) - 8
    AND OPCENTER LIKE '%s'
    GROUP BY CAST(OPCOMPDATE AS DATE)
    ORDER BY DATE;
	`, dept)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.JobCount,
			&temp.Date,
		)

		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}

//
//third party
//
func getThirdParty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]

	var temp Part
	var tempList []Part

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprintf(`
			SELECT DISTINCT TOP 20
				RUNREF,
				RUNRTNUM, 
				RUNNO,
				RUNQTY,
				SOCUST,
				--SOPO, 
				--RASOITEM, 
				ISNULL(AGPMCOMMENTS, '') AGPMCOMMENTS,
				RUNPRIORITY, 
				OPSCHEDDATE,
				ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), 0) DTDIFF,
				WCNDESC
				FROM RunsTable
				INNER JOIN RnopTable ON OPREF=RUNREF AND OPRUN= RUNNO AND RUNOPCUR=OPNO
				INNER JOIN PartTable ON PARTREF=RUNREF 
				INNER JOIN RnalTable ON RUNREF=RAREF AND RUNNO=RARUN
				INNER JOIN SohdTable ON SONUMBER=RASO
				INNER JOIN WcntTable ON OPCENTER = WCNREF
				LEFT OUTER JOIN AgcmTable ON AGPART=RUNRTNUM AND AGRUN=RUNNO
				LEFT OUTER JOIN SoitTable ON ITPART=PARTREF AND ITSO=RASO

				WHERE OPCENTER LIKE '%s'
				AND OPCOMPLETE = 0
				AND SOTYPE = 'F'
				ORDER BY RUNPRIORITY, OPSCHEDDATE`, dept)
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
			&temp.Customer,
			&temp.Comments,
			&temp.Priority,
			&temp.SchedDate,
			&temp.QueueDiff,
			&temp.WCName,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}
