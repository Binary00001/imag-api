package main

type Allocations struct {
	PartRef  string  `json:"Part_Ref"`
	Run      string  `json:"Run"`
	SO       string  `json:"SO"`
	PO       string  `json:"PO"`
	Item     string  `json:"Item"`
	Rev      string  `json:"Rev"`
	Quantity float32 `json:"Qty"`
	Type     string  `json:"Type"`
	Customer string  `json:"Customer"`
	CustDate string  `json:"Cust_Date"`
}

type CurrentLogin struct {
	EmployeeRef int    `json:"Emp_Num"`
	FirstName   string `json:"First_Name"`
	LastName    string `json:"Last_Name"`
	PartRef     string `json:"Part_Ref"`
	Run         int    `json:"Run"`
	OP          int    `json:"OP"`
	WCNum       string `json:"WC_Num"`
	PartNum     string `json:"Part_Num"`
	Description string `json:"Description"`
	WCName      string `json:"WC_Name"`
	TimeOn      int    `json:"Time_On"`
}

type Job struct {
	PartRef   string  `json:"Part_Ref"`
	PartNum   string  `json:"Part_Num"`
	Run       string  `json:"Run"`
	QueueDiff int     `json:"Queue_Diff"`
	Quantity  string  `json:"Qty"`
	Customer  *string `json:"Customer"`
	CustDate  *string `json:"Cust_Date"`
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

type Part struct {
	PartRef   string  `json:"Part_Ref"`
	PartNum   string  `json:"Part_Num"`
	Run       string  `json:"Run"`
	Quantity  float32 `json:"Qty"`
	Customer  *string `json:"Customer"`
	Comments  string  `json:"Comments"`
	Priority  int     `json:"Priority"`
	SchedDate *string `json:"Cust_Date"`
	QueueDiff int     `json:"Queue_Diff"`
	WCName    string  `json:"WC_Name"`
}

type ChartData struct {
	JobCount int    `json:"Job_Count"`
	Date     string `json:"Date"`
}

type Employee struct {
	JobCount int    `json:"Job_Count"`
	Employee string `json:"Employee"`
}

type Inventory struct {
	Sched_Date  *string  `json:"Sched_Date"`
	Req_Date    *string  `json:"Req_Date"`
	Cust        *string  `json:"Cust"`
	SO          *string  `json:"SO"`
	Item        *string  `json:"Item"`
	LTR         *string  `json:"LTR"`
	Part_Number *string  `json:"Part_Number"`
	Qty         *float32 `json:"Int"`
	Up          *string  `json:"Up"`
	On_Hand     *float32 `json:"On_Hand"`
	Lot         *string  `json:"Lot"`
	Loc         *string  `json:"Loc"`
	PO          *string  `json:"PO"`
}

type InventoryList []Inventory