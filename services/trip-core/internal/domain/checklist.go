package domain

import (
	"math"
	"time"
)

type ChecklistSection int32

const (
	ChecklistSectionUnspecified ChecklistSection = 0
	ChecklistSectionPreFlight   ChecklistSection = 1
	ChecklistSectionPacking     ChecklistSection = 2
	ChecklistSectionBuyOnSite   ChecklistSection = 3
)

type ChecklistItem struct {
	ID        string
	TripID    string
	Section   ChecklistSection
	Title     string
	IsDone    bool
	Position  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChecklistSummary struct {
	TotalItems          int32
	CompletedItems      int32
	ReadinessPercentage float64
}

func CalculateChecklistSummary(items []ChecklistItem) ChecklistSummary {
	total := int32(len(items))
	if total == 0 {
		return ChecklistSummary{
			TotalItems:          0,
			CompletedItems:      0,
			ReadinessPercentage: 0,
		}
	}

	var completed int32
	for _, it := range items {
		if it.IsDone {
			completed++
		}
	}

	pct := (float64(completed) / float64(total)) * 100.0
	pctRounded := math.Round(pct*10) / 10 // 1 decimal place

	return ChecklistSummary{
		TotalItems:          total,
		CompletedItems:      completed,
		ReadinessPercentage: pctRounded,
	}
}