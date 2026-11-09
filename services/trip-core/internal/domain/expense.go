package domain

import (
	"fmt"
	"math"
	"time"
)

type ExpenseCategory int32

const (
	ExpenseCategoryUnspecified ExpenseCategory = 0
	ExpenseCategoryStay        ExpenseCategory = 1
	ExpenseCategoryFood        ExpenseCategory = 2
	ExpenseCategoryTransport   ExpenseCategory = 3
	ExpenseCategoryTickets     ExpenseCategory = 4
	ExpenseCategoryOther       ExpenseCategory = 5
)

type Expense struct {
	ID               string
	TripID           string
	DayID            *string
	PlaceItemID      *string
	Category         ExpenseCategory
	Amount           float64
	Currency         string
	OriginalCurrency string
	OriginalAmount   float64
	FxRate           float64
	BaseAmount       float64
	IsCustomRate     bool
	Title            string
	Notes            string
	CreatedAt        time.Time
}

type BudgetSummary struct {
	BudgetLimit     float64
	TotalSpent      float64
	Currency        string
	SpentByCategory map[string]float64
	SpentByDay      map[string]float64
}

func CalculateBaseAmount(origAmount, fxRate float64) float64 {
	if fxRate <= 0 {
		fxRate = 1.0
	}
	return math.Round((origAmount*fxRate)*100) / 100
}

func CalculateBudgetSummary(limit float64, currency string, expenses []Expense) BudgetSummary {
	var total float64
	byCategory := make(map[string]float64)
	byDay := make(map[string]float64)

	for _, e := range expenses {
		effectiveAmount := e.BaseAmount
		if effectiveAmount == 0 && e.Amount > 0 {
			effectiveAmount = e.Amount
		}
		total += effectiveAmount
		catKey := fmt.Sprintf("%d", e.Category)
		byCategory[catKey] += effectiveAmount

		dayKey := "unassigned"
		if e.DayID != nil && *e.DayID != "" {
			dayKey = *e.DayID
		}
		byDay[dayKey] += effectiveAmount
	}

	// Round to 2 decimal places
	round := func(v float64) float64 {
		return math.Round(v*100) / 100
	}

	total = round(total)
	for k, v := range byCategory {
		byCategory[k] = round(v)
	}
	for k, v := range byDay {
		byDay[k] = round(v)
	}

	return BudgetSummary{
		BudgetLimit:     limit,
		TotalSpent:      total,
		Currency:        currency,
		SpentByCategory: byCategory,
		SpentByDay:      byDay,
	}
}