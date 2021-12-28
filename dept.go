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

	var temp Job
	var tempList []Job

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	query := fmt.Sprintf(`
			SELECT DISTINCT TOP 20 
      RUNREF,
      RUNRTNUM Part_Num,
      RUNNO Run,
      ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE FROM RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), 0) Queue_Diff,
      RUNQTY Qty,
      SOCUST Customer,
      CAST(ITCUSTREQ AS DATE)Cust_Date,
      RUNPRIORITY Priority,
      WCNDESC WC_Name,
      ISNULL(AGPMCOMMENTS, '') Comments,
      OPCENTER, 
      SOPO PO,
      ITNUMBER Item
 
    FROM RunsTable
      INNER JOIN RnalTable ON RUNNO=RARUN AND RUNREF=RAREF
      RIGHT OUTER JOIN PartTable ON PARTREF=RUNREF
      INNER JOIN SoitTable ON RUNREF=ITPART AND ITSO = RASO AND ITNUMBER=RASOITEM AND RASOREV=ITREV
      LEFT OUTER JOIN MrplTable ON ITSO=MRP_SONUM AND ITNUMBER=MRP_SOITEM AND ITREV=MRP_SOREV
      LEFT OUTER JOIN SohdTable ON ITSO=SONUMBER AND PALEVEL=MRP_PARTLEVEL
      LEFT JOIN RnopTable ON OPREF=RUNREF AND OPRUN=RUNNO AND OPNO=RUNOPCUR
      LEFT OUTER JOIN AgcmTable ON AGPART=RUNRTNUM AND AGRUN=RUNNO AND AGPO=SOPO AND AGITEM=ITNUMBER
      INNER JOIN WcntTable ON WCNREF=OPCENTER


    WHERE
      OPCENTER LIKE '%s' 
      AND OPCOMPDATE IS NULL AND RUNSTATUS <> 'CA' AND RUNSTATUS IS NOT NULL
    ORDER BY RUNPRIORITY, Cust_Date`, dept)
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
			&temp.QueueDiff,
			&temp.Quantity,
			&temp.Customer,
			&temp.CustDate,
			&temp.Priority,
			&temp.WCName,
			&temp.Comments,
			&temp.WCNum,
			&temp.PO,
			&temp.Item,
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
