package seinfeld

import (
	"math"
	"time"
)

type BinaryCalendar struct {
	marked map[int16][]uint32
}

func NewBinaryCalendar() Calendar {
	return &BinaryCalendar{make(map[int16][]uint32, 0)}
}

func (c *BinaryCalendar) Mark(day time.Time) {
	y, m := int16(day.Year()), int(day.Month())-1

	if _, ok := c.marked[y]; !ok {
		c.marked[y] = make([]uint32, 12)
	}

	c.marked[y][m] |= 1 << uint32(day.Day()-1)
}

func (c *BinaryCalendar) Unmark(day time.Time) {
	y, m := int16(day.Year()), int(day.Month())-1

	if _, ok := c.marked[y]; !ok {
		return
	}

	c.marked[y][m] &= c.marked[y][m] ^ (1 << uint32(day.Day()-1))
}

func (c *BinaryCalendar) IsMarked(day time.Time) bool {
	if year, ok := c.marked[int16(day.Year())]; ok {
		bit := uint32(1) << uint32(day.Day()-1)

		return year[int(day.Month())-1]&bit == bit
	}

	return false
}

func (c *BinaryCalendar) GetYear(year int) (days []time.Time) {
	y := c.marked[int16(year)]

	if y == nil {
		return []time.Time{}
	}

	result := make([]time.Time, 0)
	for m := 0; m < 12; m++ {
		result = append(result, c.formatDays(year, time.Month(m+1), y[m])...)
	}

	return result
}

func (c *BinaryCalendar) GetMonth(year int, month time.Month) []time.Time {
	y, m := c.marked[int16(year)], int(month)-1

	if y == nil {
		return []time.Time{}
	}

	return c.formatDays(year, month, y[m])
}

// Takes uint32 which represents a binary month and returns an array of dates
func (c *BinaryCalendar) formatDays(year int, month time.Month, days uint32) []time.Time {
	result := make([]time.Time, 0)

	if days == uint32(0) {
		return result
	}

	// Logarithm of the integer will tell us the highest ON bit
	for i := uint32(0); i <= uint32(math.Floor(math.Log2(float64(days)))); i++ {
		if days&(1<<i) == 1<<i {
			result = append(result, time.Date(year, month, int(i)+1, 0, 0, 0, 0, time.UTC))
		}
	}

	return result
}
