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

	tsql := fmt.Sprintf("SELECT LOTNUMBER, LOTUSERLOTID, LOTPARTREF, LOTORIGINALQTY, LOTREMAININGQTY, LOTMORUNNO, LOTLOCATION FROM LohdTable JOIN LoitTable ON LOTNUMBER=LOINUMBER WHERE LOIRECORD=(SELECT MAX(a.LOIRECORD) FROM LoitTable a WHERE LOTNUMBER=a.LOINUMBER) AND LOIPARTREF LIKE '%s' AND LOIADATE BETWEEN CAST('1995-01-01' as date) AND CAST('2021-11-13' as date)", pcmLoc)

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
