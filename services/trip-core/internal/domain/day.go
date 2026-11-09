package domain

import "time"

type TripDay struct {
	ID       string
	TripID   string
	DayIndex int32
	Date     time.Time
	Title    string
	Notes    string
}