package main

import (
	"context"
	"fmt"
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

