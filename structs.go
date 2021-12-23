package main

type Job struct {
	PartRef   string  `json:"Part_Ref"`
	PartNum   string  `json:"Part_Num"`
	Run       string  `json:"Run"`
	QueueDiff int     `json:"Queue_Diff"`
	Quantity  string  `json:"Qty"`
	Customer  *string `json:"Customer"`
	CustDate  string  `json:"Cust_Date"`
	Priority  int     `json:"Priority"`
	WCName    string  `json:"WC_Name"`
	Comments  *string `json:"Comments"`
	WCNum     string  `json:"WC_Num"`
	PO        *string `json:"PO"`
	Item      string  `json:"Item"`
}

type Burndown struct {
	PartNum   string  `json:"Part_Num"`
	Run       string  `json:"Run"`
	Comments  *string `json:"Comments"`
	WCNum     string  `json:"WC_Num"`
	WCName    string  `json:"WC_Name"`
	Quantity  string  `json:"Qty"`
	WCNNUM    string  `json:"WCNNUM"`
	QueueDiff int     `json:"Queue_Diff"`
}

type DailyGoal struct {
	Goal int `json:"Goal"`
}

type CompletedJobs struct {
	JobCount int `json:"Job_Count"`
}

type CompletedParts struct {
	PartCount int `json:"Part_Count"`
}

type DeptStats struct {
	Goal      int `json:"Goal"`
	JobCount  int `json:"Job_Count"`
	PartCount int `json:"Part_Count"`
}