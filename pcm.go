package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type pcm struct {
	PCM_ID				int			`json:"PCM_ID"`
	PCM_Location 	string  `json:"PCM_Location"`
	PCM_Num      	string  `json:"PCM_Num"`
	PCM_Sheet    	*string `json:"PCM_Sheet"`
	PCM_Rev      	*string  `json:"PCM_Rev"`
	PCM_Comments 	*string  `json:"PCM_Comments"`
}

type pcmList []pcm

func getPCMList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var temp pcm
	var tempList pcmList

	ctx := context.Background()
	
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf("SELECT * FROM PCM_LOC")

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error retrieving data: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&temp.PCM_ID,
			&temp.PCM_Location, 
			&temp.PCM_Num,
			&temp.PCM_Sheet,
			&temp.PCM_Rev,
			&temp.PCM_Comments,
		)

		if err != nil {
			log.Fatal("Error getting data: ", err.Error())
		}
		tempList = append(tempList, temp)
	}
	json.NewEncoder(w).Encode(tempList)
}