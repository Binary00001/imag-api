package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func dailyGoal(dept string) int {

	var temp DailyGoal

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
      ISNULL((SELECT MAX(DTOTAL) daily_goal
		FROM (
			SELECT OPCENTER,
			ROW_NUMBER() OVER (PARTITION BY OPCENTER ORDER BY OPCENTER) DTOTAL 
			FROM RnopTable 
			WHERE OPCOMPLETE = 0 AND OPSCHEDDATE <= CAST(GETDATE() AS DATETIME) + 30 AND OPCENTER = '%s'
		)b GROUP BY OPCENTER), 0) daily_goal`, dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}
	
	defer rows.Close()

	for rows.Next() {
			err := rows.Scan(
			&temp.Goal,
		)

		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return temp.Goal
}


func completedJobs(dept string) int {
	
	var temp CompletedJobs

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
      ISNULL((SELECT MAX(OPTOTAL)
	      FROM (
		      SELECT 
            OPCENTER, 
		        ROW_NUMBER() OVER (PARTITION BY OPCENTER ORDER BY OPCENTER) OPTOTAL 
          FROM RnopTable 
          WHERE OPCOMPDATE >= CAST(GETDATE() AS DATE) AND OPCENTER = '%s'
	      )a 
      GROUP BY OPCENTER), 0) completed_jobs`, dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}
	
	defer rows.Close()

	for rows.Next() {
			err := rows.Scan(
			&temp.JobCount,
		)

		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return temp.JobCount
}

func completedParts(dept string) int {

	var temp CompletedParts

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		 SELECT
    		CAST(ISNULL(SUM(OPACCEPT), 0) AS INT) as daily_parts  
    	FROM RnopTable 
    	WHERE OPCENTER like '%s' AND OPCOMPDATE >= CAST(GETDATE() AS DATE)`, dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}
	
	defer rows.Close()

	for rows.Next() {
			err := rows.Scan(
			&temp.PartCount,
		)

		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return temp.PartCount
}

//
//PARTS PER EMPLOYEE/DAY
//
func getEmployeeStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var temp Employee
	var tempList []Employee

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT COUNT(OPREF) AS JOBS_COMPLETED, OPINSP AS EMPLOYEE FROM RnopTable INNER JOIN RunsTable ON RUNREF = OPREF AND RunsTable.RUNNO = OPRUN WHERE RUNPKPURGED = 0 AND OPCOMPDATE >= CAST(GETDATE() AS DATE) GROUP BY OPINSP")

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Uh oh: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.JobCount,
			&temp.Employee,
		)

		if err != nil {
			log.Fatal("Error getting info: ", err.Error())
		}

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)

}

func getCurrentLogins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dept := r.URL.Query().Get("dept")
	
	var temp CurrentLogin
	var tempList []CurrentLogin

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
			ISEMPLOYEE, 
			PREMLSTNAME, 
			PREMFSTNAME, 
			ISMO, 
			ISRUN, 
			ISOP, 
			ISWCNT, 
			PARTNUM, 
			PADESC, 
			WCNDESC,
			DATEDIFF(MINUTE, ISMOSTART, GETDATE()) TIME
		FROM IstcTable
			LEFT OUTER JOIN PartTable ON PARTREF=ISMO
			LEFT JOIN WcntTable ON ISWCNT=WCNREF
			LEFT JOIN EmplTable ON ISEMPLOYEE = PREMNUMBER
		WHERE ISINDIRECT = 0 AND ISWCNT LIKE '%s'
	`, dept)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Uh oh: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.EmployeeRef,
			&temp.FirstName,
			&temp.LastName,
			&temp.PartRef,
			&temp.Run,
			&temp.OP,
			&temp.WCNum,
			&temp.PartNum,
			&temp.Description,
			&temp.WCName,
			&temp.TimeOn,
		)

		if err != nil {
			log.Fatal("Error getting info: ", err.Error())
		}

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)

}

