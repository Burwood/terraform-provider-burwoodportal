package burwoodportal

type Group struct {
	GroupName   string `json:"groupname"`
	GroupID     string  `json:"groupid"`
	Departments []Department `json:"departments"`
}

type Department struct {
	DepartmentName string `json:"departmentname"`
	DepartmentID   string  `json:"departmentid"`
	Projects    []Project `json:"projects"`
}

// type Project struct {
// 	ProjectID string `json:"projectid"`
// }

type Project struct {
	ProjectID 				string  `json:"projectid"`
	ProjectName				string	`json:"projectname"`
	PrimaryContactEmail		string	`json:"primarycontactemail"`
	BillingContactEmail		string	`json:"billingcontactemail"`
	AfterCredits 			string  `json:"aftercredits"`
	AfterCreditsAccount		string  `json:"aftercreditsaccount"`
	AfterCreditsPO			string	`json:"aftercreditspo"`
	PaidBillingAccount  	string 	`json:"paidbillingaccount"`
	TotalBudget				string	`json:"totalbudget"`
	RecurringBudget  		bool  	`json:"recurringbudget"`
	DepartmentID  			string  `json:"departmentid"`
	DepartmentName  		string  `json:"departmentname"`
}

type Allowance struct {
	PONumber       			string	 `json:"ponumber"`	
	Grant					string 	 `json:"grant"`
	Amount					int		 `json:"amount"`
	BillingAccountID		string	 `json:"billingaccountid"`
	ExpirationDate			string	 `json:"expirationdate"`
	DateSuspended			string   `json:"datesuspended"`
	DateActivated			string   `json:"dateactivated"`
	DateIssued				string	 `json:"dateissued"`
	State					string   `json:"state"`
	Recurring				bool  `json:"recurring"`
	ActualSpend				float64 `json:"actualspend"`

}

type ReportingProject struct {
	ProjectID       string  `json:"project_id"`
	CostTotal       float64 `json:"cost_total"`
	StrideDiscount  string  `json:"stride_discount"`
	I2Discount      string  `json:"i2_discount"`
	ContractCost    string  `json:"contract_cost"`
	DiscountTotal   string  `json:"discount_total"`
	Consumption     string  `json:"consumption"`
	GcpInvoiceCost  string  `json:"gcp_invoice_cost"`
	GeneralDiscount string  `json:"general_discount"`
	Markup          string  `json:"markup"`
	Adjustments     string  `json:"adjustments"`
	Subtotal        string  `json:"subtotal"`
}
