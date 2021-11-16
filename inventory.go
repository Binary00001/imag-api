package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type lot struct {
	// PCM_ID				int			`json:"PCM_ID"`
	Lot_Number 				*string  `json:"Lot_Number"`
	Lot_Id      			*string `json:"Lot_Id"`
	Lot_Part_Ref			*string `json:"Lot_Part_Ref"`
	Lot_Origin_Qty		*string	`json:"Lot_Origin_Qty"`
	Lot_Remain_Qty    *string `json:"Lot_Remain_Qty"`
	Lot_MO_Run      	*string `json:"Lot_MO_Run"`
	Lot_Loc 					*string `json:"Lot_Loc"`
}

type locationList []lot

type Inventory struct {
	Sched_Date				string	`json:"Sched_Date"`
	Req_Date					string	`json:"Req_Date"`
	Cust 							string	`json:"Cust"`
	SO 								string	`json:"SO"`
	Item							string	`json:"Item"`
	LTR								*string	`json:"LTR"`
	Part_Number				string	`json:"Part_Number"`
	Qty 							float32			`json:"Int"`
	Up 								string	`json:"Up"`
	On_Hand						float32			`json:"On_Hand"`
	Lot 							string	`json:"Lot"`
	Loc 							string	`json:"Loc"`
	PO 								string	`json:"PO"`
}

type InventoryList []Inventory

func lotsByPartNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pcmLoc := mux.Vars(r)["partNum"]
	var temp lot
	var tempList locationList

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT LOTNUMBER, RTRIM(LOTUSERLOTID) AS LOTUSERLOTID, RTRIM(LOTPARTREF) AS LOTPARTREF, LOTORIGINALQTY, LOTREMAININGQTY, LOTMORUNNO, LOTLOCATION FROM LohdTable JOIN LoitTable ON LOTNUMBER=LOINUMBER WHERE LOIRECORD=(SELECT MAX(a.LOIRECORD) FROM LoitTable a WHERE LOTNUMBER=a.LOINUMBER) AND LOIPARTREF LIKE '%s' AND LOIADATE BETWEEN CAST('1995-01-01' as date) AND CAST('2021-11-13' as date)", pcmLoc)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error retrieving data: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			// &temp.PCM_ID,
			&temp.Lot_Number,
			&temp.Lot_Id,
			&temp.Lot_Part_Ref,
			&temp.Lot_Origin_Qty,
			&temp.Lot_Remain_Qty,
			&temp.Lot_MO_Run,
			&temp.Lot_Loc,
		)

		if err != nil {
			log.Fatal("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}

func availableShip(w http.ResponseWriter, r *http.Request) {
	var temp Inventory
	var tempList InventoryList

	ctx := context.Background()
	 err := db.PingContext(ctx)
	 if err != nil {
		 log.Fatal("Could not establish a connection: ", err.Error())
	 }

	tsql := fmt.Sprintf("SELECT TOP (100) PERCENT SoitTable.ITSCHED AS SCHED, SoitTable.ITCUSTREQ AS REQD, SohdTable.SOCUST AS CUST, CONCAT(RTRIM(SOTYPE), SohdTable.SONUMBER) AS SO, SoitTable.ITNUMBER AS ITEM, SoitTable.ITREV AS LTR, RTRIM(PartTable.PARTNUM) AS PART, SoitTable.ITQTY AS QTY, SoitTable.ITDOLLARS AS UP, LohdTable.LOTREMAININGQTY AS QOH, LohdTable.LOTUSERLOTID AS LOT, LohdTable.LOTLOCATION AS LOC, SohdTable.SOPO FROM LohdTable INNER JOIN SoitTable INNER JOIN SohdTable ON SoitTable.ITSO = SohdTable.SONUMBER ON LohdTable.LOTPARTREF = SoitTable.ITPART INNER JOIN PartTable on SoitTable.ITPART = PartTable.PARTREF WHERE (LohdTable.LOTREMAININGQTY > 0) AND (SoitTable.ITACTUAL IS NULL) AND (SoitTable.ITCANCELDATE IS NULL) AND (SoitTable.ITPSITEM = 0) ORDER BY REQD, CUST")

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error retrieving data: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.Sched_Date,
			&temp.Req_Date,
			&temp.Cust,
			&temp.SO,
			&temp.Item,
			&temp.LTR,
			&temp.Part_Number,
			&temp.Qty,
			&temp.Up,
			&temp.On_Hand,
			&temp.Lot,
			&temp.Loc,
			&temp.PO,
		)

		if err != nil {
			log.Fatal("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}