package models

type CalculateRequest struct {
	BeginDate         string  `json:"beginDate"`
	PreferCreditLimit float64 `json:"preferCreditLimit"`
	PreferTenor       float64 `json:"preferTenor"`
	InterestRate      float64 `json:"interestRate"`
}

type TermSheetDetail struct {
	TermNo         int    `json:"termNo"`
	DueDate        string `json:"dueDate"`
	Installment    string `json:"installment"`
	InterestAmount string `json:"interestAmount"`
	Principal      string `json:"principal"`
	Remaining      string `json:"remaining"`
	InterestDay    int    `json:"interestDay"`
}
