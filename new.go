package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Job struct {
	PartRef					string		`json:"Part_Ref"`
	PartNum 				string		`json:"Part_Num"`
	Run 						string		`json:"Run"`
	Status 					string		`json:"Status"`
	Quantity				string 		`json:"Qty"`
	CurrentOP 			int 			`json:"Current_OP"`
	PrevOP 					int 			`json:"Prev_OP"`
	Priority				int 			`json:"Priority"`
	QueueDate				string 		`json:"Queue_Date"`
	SchedDate 			string 		`json:"Sched_Date"`
	QueueDiff				int 			`json:"Queue_Diff"`
	Comments 				string 		`json:"Comments"`
	WCName 					string 		`json:"WC_Name"`
	WCNum 					string 		`json:"WC_Num"`
}

type JobList []Job 

func getQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]

	var temp Job 
	var tempList JobList 

	ctx := context.Background()

	err := db.PingContext(ctx) 
	if err != nil {
		log.Fatal("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprintf(`
		SELECT TOP (20) 
  		RunsTable.RUNREF,
  		RunsTable.RUNRTNUM,
  		RunsTable.RUNNO,
  		RunsTable.RUNSTATUS,
  		RunsTable.RUNQTY,
  		RunsTable.RUNOPCUR,
  		TempTable.PREVOPNO,
  		RunsTable.RUNPRIORITY,
  		RnopTable.OPQDATE,
  		RnopTable.OPSCHEDDATE,
  		DATEDIFF(minute, TempTable.OPCOMPDATE, GETDATE()) QUEUETIMEMIN,
  		(SELECT AGPMCOMMENTS FROM AgcmTable WHERE A)) COMMENTS,
  		WcntTable.WCNDESC,
  		WcntTable.WCNNUM
		FROM RunsTable, RnopTable
		INNER JOIN WcntTable on WcntTable.WCNREF = RnopTable.OPCENTER,
			(SELECT RunsTable.RUNREF, RunsTable.RUNNO, RnopTable.OPNO PREVOPNO, RnopTable.OPQDATE, RnopTable.OPCOMPDATE, ROW_NUMBER() OVER (PARTITION BY RnopTable.OPREF, RnopTable.OPRUN 
				ORDER BY RnopTable.OPREF DESC,
				RnopTable.OPRUN DESC,
				RnopTable.OPNO DESC) AS RN
			FROM RunsTable, RnopTable 
			WHERE RunsTable.RUNREF = RnopTable.OPREF AND RunsTable.RUNNO = RnopTable.OPRUN AND RunsTable.RUNOPCUR > RnopTable.OPNO) AS TempTable
		WHERE
			RnopTable.OPSCHEDDATE between '2021-01-01' and CAST(GETDATE() as date)
			and RunsTable.RUNREF = TempTable.RUNREF and RunsTable.RUNNO = TempTable.RUNNO
			and RnopTable.OPREF = RunsTable.RUNREF and RnopTable.OPRUN = RunsTable.RUNNO and RnopTable.OPNO = RunsTable.RUNOPCUR
			and RnopTable.OPSHOP = '01' and RnopTable.OPCOMPLETE = 0
			and RnopTable.OPCENTER  '%s'
			and TempTable.RN = 1
		order by RunsTable.RUNPRIORITY asc
	`, dept)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal("Error executing query: ", err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&temp.PartRef,
			&temp.PartNum,
			&temp.Run,
			&temp.Status,
			&temp.Quantity,
			&temp.CurrentOP,
			&temp.PrevOP,
			&temp.Priority,
			&temp.QueueDate,
			&temp.SchedDate,
			&temp.QueueDiff,
			&temp.Comments,
			&temp.WCName,
			&temp.WCNum,
		)
		if err != nil {
			log.Fatal("Error getting data: ", err.Error())
		}
		defer rows.Close()
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}