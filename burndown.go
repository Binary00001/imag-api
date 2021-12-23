package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//
//BURNDOWN ALL
//
func getBurndown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	var temp Burndown
	var tempList []Burndown

	bd := "%burndown%"

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
				AGPART Part_Num, 
				AGRUN Run, 
				AGPMCOMMENTS Comments, 
				OPCENTER, 
				WCNDESC WC_Name, 
				RUNQTY Qty, 
				WCNNUM,
				ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), '') Queue_Diff

			FROM AgcmTable
				INNER JOIN RunsTable ON RUNRTNUM = AGPART and runno = AGRUN
				INNER JOIN RnopTable ON RUNREF = OPREF and RUNNO = oprun and RUNOPCUR = OPNO
				INNER JOIN WcntTable ON OPCENTER = WCNREF
				WHERE AGPMCOMMENTS LIKE '%s' AND 
				((RUNSTATUS <> 'CO' AND RUNSTATUS <> 'CL' AND runstatus <> 'CA') and runstatus is not null)
			ORDER BY OPCENTER ASC	`, bd)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.PartNum,
			&temp.Run,
			&temp.Comments,
			&temp.WCNum,
			&temp.WCName,
			&temp.Quantity,
			&temp.WCNNUM,
			&temp.QueueDiff,
		)

		if err != nil {
			fmt.Println("Error: ", err.Error())
		}


		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}

//
//BURNDOWN BY DEPT
//
func getDeptBurndown(w http.ResponseWriter, r *http.Request) {
	dept := mux.Vars(r)["dept"]
	bd := "%BURNDOWN%"
	w.Header().Set("Content-Type", "application/json")

	var temp Burndown
	var tempList []Burndown

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
				AGPART Part_Num, 
				AGRUN Run, 
				AGPMCOMMENTS Comments, 
				OPCENTER, 
				WCNDESC WC_Name, 
				RUNQTY Qty, 
				WCNNUM,
				ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), '') Queue_Diff

			FROM AgcmTable
				INNER JOIN RunsTable ON RUNRTNUM = AGPART and runno = AGRUN
				INNER JOIN RnopTable ON RUNREF = OPREF and RUNNO = oprun and RUNOPCUR = OPNO
				INNER JOIN WcntTable ON OPCENTER = WCNREF
				WHERE AGPMCOMMENTS LIKE '%s' AND OPCENTER = '%s' AND
				((RUNSTATUS <> 'CO' AND RUNSTATUS <> 'CL' AND runstatus <> 'CA') and runstatus is not null)
			ORDER BY OPCENTER ASC`, bd, dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}
	
	defer rows.Close()

	for rows.Next() {
			err := rows.Scan(
			&temp.PartNum,
			&temp.Run,
			&temp.Comments,
			&temp.WCNum,
			&temp.WCName,
			&temp.Quantity,
			&temp.WCNNUM,
			&temp.QueueDiff,
		)

		if err != nil {
			fmt.Println("Error: ", err.Error())
		}


		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}