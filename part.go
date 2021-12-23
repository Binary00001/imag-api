package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getParts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vals := r.URL.Query()
	po := vals.Get("po")
	line := vals.Get("line")
	run := vals.Get("run")
	part := vals.Get("part")

	var parts []Job

	tsql := fmt.Sprintf(`
			SELECT DISTINCT TOP 20 
			RUNREF, 
			RUNRTNUM, 
			RUNNO, 
			RUNQTY, 
			AGPMCOMMENTS, 
		CASE
			WHEN RTRIM(SOCUST) = 'BOECOM' THEN ('BOE' + SUBSTRING(SOPO,0,4))
			ELSE 
			SOCUST
		END CUSTOMER,
		SOPO,
		ITNUMBER, 
		CAST(ITCUSTREQ AS DATE)CUSTDATE,
		OPCENTER, 
		(SELECT WCNDESC FROM WcntTable WHERE WCNNUM = OPCENTER) WCDESC,
		RUNPRIORITY,
		ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), '') QUEUETIME
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
		LEFT OUTER JOIN AgcmTable ON AGPART = RUNRTNUM AND AGRUN = RUNNO
			AND AGPO = SOPO AND AGITEM = ITNUMBER	

		WHERE SOPO = '%s' AND LTRIM(ITNUMBER) = '%s' AND RUNNO = '%s' AND RUNRTNUM = '%s'
		AND ITCANCELED = 0 
		AND ITINVOICE = 0 
		AND ITPSSHIPPED = 0 
		AND ((RUNSTATUS <> 'CO' AND RUNSTATUS <> 'CA' AND RUNSTATUS <> 'CL') OR RUNSTATUS IS NULL)
		AND OPCOMPLETE = 0 
		AND RUNOPCUR = OPNO
		ORDER BY RUNPRIORITY, CUSTDATE ASC`, po, line, run, part)

	parts = getJobs(tsql)
	// 	if err != nil {
	// 		fmt.Println("Error getting data: ", err.Error())
	// 	}
	// 	tempList = append(tempList, temp)
	// }
	json.NewEncoder(w).Encode(parts)
}