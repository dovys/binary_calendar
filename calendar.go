package seinfeld

import "time"

type Calendar interface {
	Mark(day time.Time)
	Unmark(day time.Time)
	IsMarked(day time.Time) bool
	GetMonth(year int, month time.Month) (days []time.Time)
	GetYear(year int) (days []time.Time)
}
