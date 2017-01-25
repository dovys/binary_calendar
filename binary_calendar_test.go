package seinfeld

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMark(t *testing.T) {
	now := time.Now()

	cases := []struct {
		test     string
		mark     time.Time
		compare  time.Time
		expected bool
	}{
		{"Simple case", now, now, true},
		{"Simple case for last year", now.AddDate(-1, 0, 0), now.AddDate(-1, 0, 0), true},
		{"Mark different day", now, now.AddDate(0, 0, -1), false},
		{"Mark different month", now, now.AddDate(0, -1, 0), false},
		{"Mark different year", now.AddDate(-2, 0, 0), now, false},
	}

	for _, test := range cases {
		t.Run(test.test, func(t *testing.T) {
			c := NewBinaryCalendar()
			c.Mark(test.mark)

			assert.Equal(t, test.expected, c.IsMarked(test.compare))
		})
	}
}

func TestMultipleMarks(t *testing.T) {
	c := NewBinaryCalendar()
	march := time.Date(2016, 3, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i <= 30; i++ {
		c.Mark(march.AddDate(0, 0, i))
	}

	for i := 0; i <= 30; i++ {
		assert.True(t, c.IsMarked(march.AddDate(0, 0, i)))
		assert.False(t, c.IsMarked(march.AddDate(-1, 0, i)))
		assert.False(t, c.IsMarked(march.AddDate(0, -2, i)))
	}
}

func TestUnmark(t *testing.T) {
	now := time.Now()
	c := NewBinaryCalendar()

	c.Mark(now.AddDate(0, 0, 1))
	c.Mark(now)
	c.Mark(now.AddDate(1, 0, 0))
	c.Mark(now.AddDate(0, 1, 0))
	c.Mark(now.AddDate(-1, -1, 0))
	c.Mark(now.AddDate(-1, -1, -1))

	c.Unmark(now.AddDate(0, 0, 1))
	c.Unmark(now)
	c.Unmark(now)
	c.Unmark(now.AddDate(1, 1, 1))
	c.Unmark(now.AddDate(-1, -1, -1))

	assert.True(t, c.IsMarked(now.AddDate(1, 0, 0)))
	assert.True(t, c.IsMarked(now.AddDate(-1, -1, 0)))
	assert.False(t, c.IsMarked(now.AddDate(-1, -1, -1)))
	assert.False(t, c.IsMarked(now.AddDate(1, 1, 1)))
	assert.False(t, c.IsMarked(now.AddDate(0, 0, 1)))
	assert.False(t, c.IsMarked(now))
}

func TestBitmasking(t *testing.T) {
	c := &BinaryCalendar{make(map[int16][]uint32, 0)}

	march := time.Date(2016, 3, 1, 12, 0, 0, 0, time.UTC)

	// Fill in March
	for i := 0; i < 31; i++ {
		c.Mark(march.AddDate(0, 0, i))
	}

	// Assert that the first 31 (0-30) bits are turned on
	assert.Equal(t, uint32(0x7fffffff), c.marked[2016][2])

	// Unmark March 15th
	c.Unmark(march.AddDate(0, 0, 14))

	// Assert that first 31 bits except the 15th are turned on
	assert.Equal(t, uint32(0x7fffBfff), c.marked[2016][2])

	// Assert that January is empty
	assert.Equal(t, uint32(0), c.marked[2016][0])
}

func TestGetMonthAndYear(t *testing.T) {
	c := NewBinaryCalendar()

	dateStrings := []string{
		"2005-02-05",
		"2006-01-01",
		"2006-02-01",
		"2006-02-15",
		"2006-02-28",
		"2006-03-01",
		"2006-11-30",
		"2006-12-31",
		"2007-01-01",
	}
	dates := make([]time.Time, len(dateStrings))

	for i, s := range dateStrings {
		d, e := time.Parse(time.RFC3339, s+"T00:00:00Z")
		assert.Nil(t, e)

		dates[i] = d
		c.Mark(d)
	}

	assert.Equal(t, dates[2:5], c.GetMonth(2006, time.February))
	assert.Equal(t, dates[5:6], c.GetMonth(2006, time.March))
	assert.Equal(t, dates[0:1], c.GetMonth(2005, time.February))
	assert.Equal(t, []time.Time{}, c.GetMonth(2006, time.July))
	assert.Equal(t, []time.Time{}, c.GetMonth(2005, time.March))

	assert.Equal(t, []time.Time{}, c.GetYear(2004))
	assert.Equal(t, dates[0:1], c.GetYear(2005))
	assert.Equal(t, dates[1:8], c.GetYear(2006))
	assert.Equal(t, dates[8:9], c.GetYear(2007))
	assert.Equal(t, []time.Time{}, c.GetYear(2008))
}

func setupLargeCalendar() Calendar {
	c := NewBinaryCalendar()
	for y := 30; y >= 0; y-- {
		for d := 360; d >= 0; d-- {
			c.Mark(time.Now().AddDate(-y, 0, -d))
		}
	}

	return c
}

func BenchmarkMarkSameDay(b *testing.B) {
	b.ReportAllocs()
	c := NewBinaryCalendar()
	now := time.Now()

	for i := 0; i < b.N; i++ {
		c.Mark(now)
	}
}

func BenchmarkMarkIncrementalDays(b *testing.B) {
	b.ReportAllocs()
	c := NewBinaryCalendar()
	now := time.Now()

	for i := 0; i < b.N; i++ {
		c.Mark(now.AddDate(0, 0, i))
	}
}

func BenchmarkMarkLargeCalendar(b *testing.B) {
	b.ReportAllocs()
	c := setupLargeCalendar()
	now := time.Now()

	for i := 0; i < b.N; i++ {
		c.Mark(now.AddDate(0, 0, i))
	}
}

func BenchmarkMarkUnmark(b *testing.B) {
	b.ReportAllocs()
	c := NewBinaryCalendar()
	now := time.Now()

	for i := 0; i < b.N; i++ {
		c.Mark(now)
		c.Unmark(now)
	}
}

func BenchmarkUnmarkLargeCalendar(b *testing.B) {
	b.ReportAllocs()
	c := setupLargeCalendar()
	now := time.Now()

	for i := 0; i < b.N; i++ {
		c.Unmark(now.AddDate(0, 0, -i))
	}
}

func BenchmarkIsMarked(b *testing.B) {
	b.ReportAllocs()
	c := setupLargeCalendar()
	now := time.Now()

	for i := 0; i < b.N; i++ {
		c.IsMarked(now.AddDate(0, 0, -i))
	}
}

func BenchmarkGetMonth(b *testing.B) {
	b.ReportAllocs()
	c := setupLargeCalendar()

	benchmarks := []struct {
		test string
		date time.Time
	}{
		{"One month back", time.Now().AddDate(0, -1, 0)},
		{"One year back", time.Now().AddDate(-1, 0, 0)},
		{"Five years back", time.Now().AddDate(-5, 0, 0)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.test, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.GetMonth(bm.date.Year(), bm.date.Month())
			}
		})
	}
}

func BenchmarkGetYear(b *testing.B) {
	b.ReportAllocs()
	c := setupLargeCalendar()

	for i := 0; i < b.N; i++ {
		c.GetYear(time.Now().AddDate(-5, 0, 0).Year())
	}
}
