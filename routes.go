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
	Part_Number string  `json:"Part_Number"`
	Run         *string `json:"Run"`
	PO          string  `json:"PO"`
	Item        *string `json:"Item"`
	Queue_Days  *int    `json:"Queue_Days"`
	Customer    string  `json:"Customer"`
	Priority    int     `json:"Priority"`
	Comments    *string `json:"Comments"`
	Due_Date    string  `json:"Due_Date"`
	Run_Qty     string  `json:"Run_Qty"`
	Work_Center string  `json:"Work_Center"`
	WC_Num      string  `json:"WC_Num"`
	WC_Name     string  `json:"WC_Name"`
}

type jobList []job

type department struct {
	Part_Ref      string  `json:"Part_Ref"`
	Run_Num       string  `json:"Run_Num"`
	Run           *string `json:"Run"`
	Work_Center   *string `json:"Work_Center"`
	WC_Name       *string `json:"WC_Name"`
	Run_Status    string  `json:"Run_Status"`
	Run_Qty       string  `json:"Run_Qty"`
	Current_OP    *string `json:"Current_OP"`
	OP_QDate      *string `json:"OP_QDate"`
	OP_SchedDate  string  `json:"OP_SchedDate"`
	Run_Priority  int     `json:"Run_Priority"`
	OP_No         int     `json:"OP_No"`
	Prev_QDate    *string `json:"Prev_QDate"`
	Prev_CompDate *string `json:"Prev_CompDate"`
	Comments      *string `json:"Comments"`
	Date_DiffQ    int     `json:"Date_DiffQ"`
	Date_DiffNow  int     `json:"Date_DiffNow"`
}

type deptList []department

type employee struct {
	Job_Count int    `json:"Job_Count"`
	Employee  string `json:"Employee"`
}

type employeeList []employee

type Burndown struct {
	Ref         string `json:"Ref"`
	Part_Number string `json:"Part_Number"`
	Run         int    `json:"Run"`
	Qty         string `json:"Qty"`
	Comments    string `json:"Comments"`
	Customer    string `json:"Customer"`
	Work_Center string `json:"Work_Center"`
	Cust_Date   string `json:"Cust_Date"`
	WC_Name     string `json:"WC_Name"`
	Queue_Time  int    `json:"Queue_Time"`
	Priority    int    `json:"Priority"`
}

type BurndownList []Burndown

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

	tsql := fmt.Sprintf("SELECT COUNT(OPREF) AS JOBS_COMPLETED, OPINSP AS EMPLOYEE FROM RnopTable INNER JOIN RunsTable ON RUNREF = OPREF AND RunsTable.RUNNO = OPRUN WHERE RUNPKPURGED = 0 AND OPCOMPDATE >= CAST(GETDATE() AS DATE) GROUP BY OPINSP")

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

	var temp department
	var tempList deptList

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish connection: ", err.Error())
	}

	// tsql := fmt.Sprintf("SELECT TOP 20 RTRIM(PART_NUMBER) AS PART_NUMBER, RTRIM(RUN) AS RUN, RTRIM(PO) AS PO, RTRIM(LTRIM(ITEM)) AS ITEM, DAYS_IN_QUEUE, RTRIM(CUSTOMER) AS CUSTOMER, PRIORITY, COMMENTS, CAST(CUST_REQ_DATE AS DATETIME) AS CUST_REQ_DATE, RUN_QTY, RTRIM(WORK_CENTER) AS WORK_CENTER, RTRIM(WC) AS WC, RTRIM(t2.WCNDESC) AS WC_NAME FROM QueueInfo INNER JOIN WcntTable AS t2 ON WC = t2.WCNNUM WHERE WC = '%s' ORDER BY CUST_REQ_DATE ASC", dept)

	// tsql := fmt.Sprintf("select distinct TOP 20 RTRIM(runstable.Runref) As RUNREF, RTRIM(runstable.runrtnum) AS RUNRTNUM, runstable.runno, RTRIM(OPCENTER) AS OPCENTER, RUNSTATUS, RUNQTY, runopcur,RnopTable.OPQDATE,  RnopTable.OPSCHEDDATE,RunsTable.RUNPRIORITY, f.PrevOPNO, f.OPQDATE PREVQDATE, f.OPCOMPDATE PREVCOMPDATE, (SELECT AGPMCOMMENTS FROM AgcmTable WHERE AGPART = RUNRTNUM AND AGRUN=RunsTable.RUNNO) AS COMMENTS, DATEDIFF(day,f.OPCOMPDATE, RnopTable.OPQDATE) DtDiffQue, DATEDIFF(day, f.OPCOMPDATE, GETDATE()) DtDiffNow from runstable WITH (nolock),RnopTable WITH (nolock), (select a.runref, a.runno, b.OPNO PrevOPNO, OPQDATE, OPCOMPDATE, ROW_NUMBER() OVER (PARTITION BY opref, oprun ORDER BY opref DESC, oprun DESC, OPNO desc) as rn from runstable a,rnopTable b where a.runref = b.opref and a.RunNO = b.oprun And b.opno < a.runopcur) as f where RnopTable.OPSCHEDDATE between '2021-01-01' and CAST(GETDATE() AS DATE) and RunsTable.runref =  f.runref AND RunsTable.runno = f.runno and RnopTable.opref = RunsTable.runref AND RunsTable.runno = RnopTable.OPRUN and RunsTable.runopcur = RnopTable.Opno and RnopTable.OPSHOP = '01' AND RnopTable.OPCENTER = '%s' AND RnopTable.OPCOMPLETE = 0 and f.rn = 1 ORDER BY RunsTable.RUNPRIORITY, RnopTable.OPSCHEDDATE", dept)
	tsql := fmt.Sprintf("SELECT TOP 20 RTRIM(RunsTable.RUNREF) As Run_Ref, RTRIM(RunsTable.RUNRTNUM) AS Part_Number, RunsTable.RUNNO AS Run, RTRIM(OPCENTER) AS WC_Ref, WcntTable.WCNDESC AS WC_Name, RunsTable.RUNSTATUS AS Run_Status, RunsTable.RUNQTY AS Qty, RunsTable.RUNOPCUR AS Current_Op, RnopTable.OPQDATE, RnopTable.OPSCHEDDATE, RunsTable.RUNPRIORITY AS Run_Priority, f.PrevOPNO, f.OPQDATE PREVQDATE, f.OPCOMPDATE Prev_Comp_Date, (SELECT AGPMCOMMENTS FROM AgcmTable WHERE AGPART = RUNRTNUM AND AGRUN=RunsTable.RUNNO) AS Comments, DATEDIFF(day,f.OPCOMPDATE, RnopTable.OPQDATE) DtDiffQue, DATEDIFF(day, f.OPCOMPDATE, GETDATE()) DtDiffNow from runstable,RnopTable INNER JOIN WcntTable ON WCNREF=OPCENTER, (select a.runref, a.runno, b.OPNO PrevOPNO, OPQDATE, OPCOMPDATE, ROW_NUMBER() OVER (PARTITION BY opref, oprun ORDER BY opref DESC, oprun DESC, OPNO desc) as rn from runstable a,rnopTable b where a.runref = b.opref and a.RunNO = b.oprun And b.opno < a.runopcur) as f where RnopTable.OPSCHEDDATE between '2021-01-01' and CAST(GETDATE() AS DATE) and RunsTable.runref =  f.runref AND RunsTable.runno = f.runno and RnopTable.opref = RunsTable.runref AND RunsTable.runno = RnopTable.OPRUN and RunsTable.runopcur = RnopTable.Opno and RnopTable.OPSHOP LIKE '01' AND RnopTable.OPCENTER = '%s' AND RnopTable.OPCOMPLETE = 0 and f.rn = 1 ORDER BY RunsTable.RUNPRIORITY, RnopTable.OPSCHEDDATE", dept)
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error executing query: ", err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&temp.Part_Ref,
			&temp.Run_Num,
			&temp.Run,
			&temp.Work_Center,
			&temp.WC_Name,
			&temp.Run_Status,
			&temp.Run_Qty,
			&temp.Current_OP,
			&temp.OP_QDate,
			&temp.OP_SchedDate,
			&temp.Run_Priority,
			&temp.OP_No,
			&temp.Prev_QDate,
			&temp.Prev_CompDate,
			&temp.Comments,
			&temp.Date_DiffQ,
			&temp.Date_DiffNow,
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

	var temp Burndown
	var tempList BurndownList

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish connection: ", err.Error())
	}

	tsql := fmt.Sprint(`
	SELECT RUNREF, RUNRTNUM, RUNNO, RUNQTY, AGPMCOMMENTS, 
	CASE
		WHEN RTRIM(SOCUST) = 'BOECOM' THEN ('BOE' + SUBSTRING(SOPO,0,4))
		ELSE 
		SOCUST
	 END CUSTOMER, 
	OPCENTER, 
	CAST(ITCUSTREQ AS DATE)CUST_DATE,
	(SELECT WCNDESC FROM WcntTable WHERE WCNNUM = OPCENTER) WCNO,
	ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), '') QUEUETIME,
	RUNPRIORITY
	FROM RunsTable
	INNER JOIN RnalTable ON RUNNO = RARUN 
		AND RUNREF = RAREF
	RIGHT OUTER JOIN PartTable
	INNER JOIN SoitTable ON PARTREF = ITPART ON RASO = ITSO 
		AND RASOITEM = ITNUMBER 
		AND RASOREV = ITREV
	LEFT OUTER JOIN MrplTable ON ITSO = MRP_SONUM 
		AND ITNUMBER = MRP_SOITEM
		AND ITREV = MRP_SOREV
	LEFT OUTER JOIN SohdTable ON ITSO = SONUMBER
		AND PALEVEL = MRP_PARTLEVEL
	INNER JOIN RnopTable ON OPREF = RUNREF AND OPRUN = RUNNO 
	INNER JOIN AgcmTable ON AGPART = RUNRTNUM AND AGRUN = RUNNO
		AND AGPO = SOPO AND AGITEM = ITNUMBER	

	WHERE ITCANCELED = 0 
	AND ITINVOICE = 0 
	AND ITPSSHIPPED = 0 
	AND AGPMCOMMENTS LIKE '%BURNDOWN%'
	AND ((RUNSTATUS <> 'CO' AND RUNSTATUS <> 'CA' AND RUNSTATUS <> 'CL') OR RUNSTATUS IS NULL)
	AND OPCOMPLETE = 0 
	AND RUNOPCUR = OPNO
	ORDER BY OPCENTER ASC`)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error executing query: ", err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&temp.Ref,
			&temp.Part_Number,
			&temp.Run,
			&temp.Qty,
			&temp.Comments,
			&temp.Customer,
			&temp.Work_Center,
			&temp.Cust_Date,
			&temp.WC_Name,
			&temp.Queue_Time,
			&temp.Priority,
		)

		if err != nil {
			log.Fatal("Error: ", err.Error())
		}

		defer rows.Close()

		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}
