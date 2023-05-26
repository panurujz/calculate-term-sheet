package services

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/panurujz/calculate-term-sheet/models"
)

var df_yyyy_mm_dd = "2006-01-02"
var df_dd_mmm_yyyy = "02-Jan-2006"

func CalculateTs(c echo.Context) (err error) {

	r := new(models.CalculateRequest)
	if err = c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	fmt.Println("################################################################")
	fmt.Println("###################### CALCULATE TERM-SHEET ####################")
	fmt.Println("################################################################")

	beginDate, _ := time.Parse("2006-01-02", r.BeginDate)
	preferCreditLimit := r.PreferCreditLimit
	preferTenor := r.PreferTenor
	interestRate := r.InterestRate
	installment := calculateInstallment(preferCreditLimit, preferTenor, interestRate)

	var remaining float64 = preferCreditLimit
	var totalDay float64
	var totalInstallment float64
	var totalInterest float64
	var totalPrincipal float64
	var termNo = 0

	fmt.Printf("beginDate = %s\n", beginDate.Format(df_yyyy_mm_dd))
	fmt.Printf("preferCreditLimit = %.2f\n", preferCreditLimit)
	fmt.Printf("preferTenor = %.2f\n", preferTenor)
	fmt.Printf("interestRate = %.2f\n", interestRate)
	fmt.Printf("installment = %.2f\n", installment)

	var details []models.TermSheetDetail

	for termNo < int(preferTenor) {
		termNo++

		dueDate := beginDate
		nextDueDate := dueDate.AddDate(0, 1, 0)
		interestDay := nextDueDate.Sub(dueDate).Hours() / 24
		interestAmount := calculateInterestAmount(remaining, interestRate, interestDay)
		principal := installment - interestAmount

		if termNo == int(preferTenor) {
			if remaining <= installment {
				installment = remaining + interestAmount
			}
			principal = remaining
			remaining = 0.0
		} else {
			remaining = remaining - principal
		}

		if termNo == 1 {
			fmt.Println("-----------------------------------------------------------------------------")
			fmt.Println("termNo|nextDueDate|installment|interestAmount|principal|remaining|interestDay")
			fmt.Println("-----------------------------------------------------------------------------")
		}

		d := models.TermSheetDetail{
			TermNo:         termNo,
			DueDate:        dueDate.Format(df_dd_mmm_yyyy),
			Installment:    fmt.Sprintf("%.2f", installment),
			InterestAmount: fmt.Sprintf("%.2f", interestAmount),
			Principal:      fmt.Sprintf("%.2f", principal),
			Remaining:      fmt.Sprintf("%.2f", remaining),
			InterestDay:    int(interestDay),
		}

		details = append(details, d)

		totalDay = totalDay + interestDay
		totalInstallment = totalInstallment + installment
		totalInterest = totalInterest + interestAmount
		totalPrincipal = totalPrincipal + principal
		beginDate = nextDueDate
	}

	for _, v := range details {
		fmt.Printf("%6d|%s|%s|%s|%s|%s|%d\n",
			v.TermNo,
			v.DueDate,
			v.Installment,
			v.InterestAmount,
			v.Principal,
			v.Remaining,
			v.InterestDay,
		)
	}

	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("Total -----> |%.2f|%.2f|%.2f|%d\n",
		totalInstallment,
		totalInterest,
		totalPrincipal,
		int(totalDay),
	)
	fmt.Println("-----------------------------------------------------------------------------")

	return c.JSON(http.StatusOK, details)
}

func calculateInstallment(preferCreditLimit, preferTenor, interestRate float64) (installment float64) {
	ratio := math.Pow(10, float64(-1))
	ratePerMonth := calculateRatePerMonth(interestRate)
	interestRateByTenor := ratePerMonth + 1.0
	taskResult := math.Pow(interestRateByTenor, preferTenor)
	taskResult = 1 / taskResult
	taskResult = (1 - taskResult) / ratePerMonth
	taskResult = preferCreditLimit / taskResult
	return math.Round(taskResult*ratio) / ratio
}

func calculateRatePerMonth(interestRate float64) (ratePerMonth float64) {
	return (interestRate / 1200.0) * 1.0
}

func calculateInterestAmount(remaining, interestRate, interestDay float64) (interestAmount float64) {
	taskResult := interestRate / 100
	taskResult = remaining * taskResult * interestDay
	return taskResult / 365.0
}
